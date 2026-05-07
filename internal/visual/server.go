package visual

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	wsMagicGUID   = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	idleTimeout   = 30 * time.Minute
	lifecycleTick = 60 * time.Second
	pollInterval  = 2 * time.Second
)

// Config holds server configuration.
type Config struct {
	Port     int
	Host     string
	URLHost  string
	ContentDir string
	StateDir   string
	OwnerPID   int
	FrameHTML  string
	HelperJS   string
}

// Server is a visual brainstorming companion server.
type Server struct {
	cfg          Config
	clients      map[net.Conn]struct{}
	clientsMu    sync.Mutex
	lastActivity time.Time
	knownFiles   map[string]bool
	knownMu      sync.Mutex
	stopCh       chan struct{}
}

// NewServer creates a new visual brainstorming server.
func NewServer(cfg Config) *Server {
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.URLHost == "" {
		if cfg.Host == "127.0.0.1" {
			cfg.URLHost = "localhost"
		} else {
			cfg.URLHost = cfg.Host
		}
	}
	return &Server{
		cfg:        cfg,
		clients:    make(map[net.Conn]struct{}),
		knownFiles: make(map[string]bool),
		stopCh:     make(chan struct{}),
	}
}

// StartInfo is the JSON output emitted on server startup.
type StartInfo struct {
	Type      string `json:"type"`
	Port      int    `json:"port"`
	Host      string `json:"host"`
	URLHost   string `json:"url_host"`
	URL       string `json:"url"`
	ScreenDir string `json:"screen_dir"`
	StateDir  string `json:"state_dir"`
}

// Start starts the server and blocks until it shuts down.
func (s *Server) Start() error {
	os.MkdirAll(s.cfg.ContentDir, 0755)
	os.MkdirAll(s.cfg.StateDir, 0755)

	// Seed known files
	entries, err := os.ReadDir(s.cfg.ContentDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".html") {
				s.knownFiles[e.Name()] = true
			}
		}
	}

	s.lastActivity = time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/files/", s.handleFile)

	srv := &http.Server{
		Handler: mux,
		ConnState: func(c net.Conn, cs http.ConnState) {
			if cs == http.StateNew {
				s.touchActivity()
			}
		},
	}

	// WebSocket upgrade listener
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			s.handleWebSocket(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	// Start lifecycle monitor
	lifecycleDone := make(chan struct{})
	go func() {
		defer close(lifecycleDone)
		ticker := time.NewTicker(lifecycleTick)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !s.ownerAlive() {
					s.shutdown(srv, "owner process exited")
					return
				}
				if time.Since(s.lastActivity) > idleTimeout {
					s.shutdown(srv, "idle timeout")
					return
				}
			case <-s.stopCh:
				return
			}
		}
	}()

	// Start file poller
	pollerDone := make(chan struct{})
	go func() {
		defer close(pollerDone)
		s.pollFiles()
	}()

	// Validate owner PID at startup
	if s.cfg.OwnerPID != 0 {
		if !s.ownerAlive() {
			log.Printf(`{"type":"owner-pid-invalid","pid":%d,"reason":"dead at startup"}`, s.cfg.OwnerPID)
			s.cfg.OwnerPID = 0
		}
	}

	// Listen
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	info := StartInfo{
		Type:      "server-started",
		Port:      ln.Addr().(*net.TCPAddr).Port,
		Host:      s.cfg.Host,
		URLHost:   s.cfg.URLHost,
		URL:       fmt.Sprintf("http://%s:%d", s.cfg.URLHost, ln.Addr().(*net.TCPAddr).Port),
		ScreenDir: s.cfg.ContentDir,
		StateDir:  s.cfg.StateDir,
	}
	infoJSON, _ := json.Marshal(info)
	log.Println(string(infoJSON))

	// Write server-info
	infoPath := filepath.Join(s.cfg.StateDir, "server-info")
	os.WriteFile(infoPath, append(infoJSON, '\n'), 0644)

	err = srv.Serve(ln)
	<-pollerDone
	<-lifecycleDone
	return err
}

// Stop signals the server to shut down.
func (s *Server) Stop() {
	close(s.stopCh)
}

func (s *Server) shutdown(srv *http.Server, reason string) {
	log.Printf(`{"type":"server-stopped","reason":"%s"}`, reason)
	infoPath := filepath.Join(s.cfg.StateDir, "server-info")
	os.Remove(infoPath)
	stoppedPath := filepath.Join(s.cfg.StateDir, "server-stopped")
	stoppedJSON, _ := json.Marshal(map[string]any{"reason": reason, "timestamp": time.Now().UnixMilli()})
	os.WriteFile(stoppedPath, append(stoppedJSON, '\n'), 0644)
	srv.Close()
}

func (s *Server) touchActivity() {
	s.lastActivity = time.Now()
}

func (s *Server) ownerAlive() bool {
	if s.cfg.OwnerPID == 0 {
		return true
	}
	proc, err := os.FindProcess(s.cfg.OwnerPID)
	if err != nil {
		return false
	}
	err = proc.Signal(os.Interrupt)
	// On Unix, Signal returns nil if process exists (permission or signal sent)
	// On Windows, OpenProcess succeeds if PID exists
	return err == nil
}

// ========== HTTP Handlers ==========

var mimeTypes = map[string]string{
	".html": "text/html; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".js":   "application/javascript; charset=utf-8",
	".json": "application/json",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
}

const waitingPage = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>rein Brainstorm</title>
<style>body{font-family:system-ui,sans-serif;padding:2rem;max-width:800px;margin:0 auto}
h1{color:#333}p{color:#666}</style>
</head>
<body><h1>rein Brainstorm</h1>
<p>Waiting for content...</p></body></html>`

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.touchActivity()
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	screenFile := s.newestScreen()
	var html string
	if screenFile != "" {
		raw, err := os.ReadFile(screenFile)
		if err != nil {
			html = waitingPage
		} else {
			content := string(raw)
			if isFullDocument(content) {
				html = content
			} else {
				html = wrapInFrame(s.cfg.FrameHTML, content)
			}
		}
	} else {
		html = waitingPage
	}

	// Inject helper script
	helperInjection := "<script>\n" + s.cfg.HelperJS + "\n</script>"
	if idx := strings.LastIndex(html, "</body>"); idx != -1 {
		html = html[:idx] + helperInjection + "\n" + html[idx:]
	} else {
		html += helperInjection
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	s.touchActivity()
	fileName := filepath.Base(r.URL.Path)
	filePath := filepath.Join(s.cfg.ContentDir, fileName)

	if _, err := os.Stat(filePath); err != nil {
		http.NotFound(w, r)
		return
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	ct, ok := mimeTypes[ext]
	if !ok {
		ct = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ct)
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write(data)
}

// ========== WebSocket Handler ==========

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		http.Error(w, "missing Sec-WebSocket-Key", http.StatusBadRequest)
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "websocket not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		return
	}

	// Compute accept key
	h := sha1.New()
	h.Write([]byte(key + wsMagicGUID))
	accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Write handshake response
	handshake := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	conn.Write([]byte(handshake))

	s.addClient(conn)

	go s.readWebSocket(conn)
}

func (s *Server) addClient(conn net.Conn) {
	s.clientsMu.Lock()
	s.clients[conn] = struct{}{}
	s.clientsMu.Unlock()
}

func (s *Server) removeClient(conn net.Conn) {
	s.clientsMu.Lock()
	delete(s.clients, conn)
	s.clientsMu.Unlock()
	conn.Close()
}

func (s *Server) readWebSocket(conn net.Conn) {
	defer s.removeClient(conn)

	buf := make([]byte, 4096)
	var pending []byte

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		pending = append(pending, buf[:n]...)

		for len(pending) > 0 {
			frame, consumed, err := decodeFrame(pending)
			if err != nil {
				s.sendClose(conn, 1002, "protocol error")
				return
			}
			if frame == nil {
				break // need more data
			}
			pending = pending[consumed:]

			switch frame.opcode {
			case 0x01: // TEXT
				s.handleMessage(string(frame.payload), conn)
			case 0x08: // CLOSE
				s.sendClose(conn, 0, "")
				return
			case 0x09: // PING
				s.writeFrame(conn, 0x0A, frame.payload)
			case 0x0A: // PONG
				// ignore
			default:
				s.sendClose(conn, 1003, "unsupported opcode")
				return
			}
		}
	}
}

type wsFrame struct {
	opcode  byte
	payload []byte
}

func decodeFrame(data []byte) (*wsFrame, int, error) {
	if len(data) < 2 {
		return nil, 0, nil
	}

	opcode := data[0] & 0x0F
	secondByte := data[1]
	masked := (secondByte & 0x80) != 0
	payloadLen := int(secondByte & 0x7F)
	offset := 2

	if !masked {
		return nil, 0, fmt.Errorf("client frames must be masked")
	}

	if payloadLen == 126 {
		if len(data) < 4 {
			return nil, 0, nil
		}
		payloadLen = int(binary.BigEndian.Uint16(data[2:4]))
		offset = 4
	} else if payloadLen == 127 {
		if len(data) < 10 {
			return nil, 0, nil
		}
		payloadLen = int(binary.BigEndian.Uint64(data[2:10]))
		offset = 10
	}

	maskStart := offset
	dataStart := offset + 4
	totalLen := dataStart + payloadLen
	if len(data) < totalLen {
		return nil, 0, nil
	}

	mask := data[maskStart:dataStart]
	payload := make([]byte, payloadLen)
	for i := 0; i < payloadLen; i++ {
		payload[i] = data[dataStart+i] ^ mask[i%4]
	}

	return &wsFrame{opcode: opcode, payload: payload}, totalLen, nil
}

func (s *Server) writeFrame(conn net.Conn, opcode byte, payload []byte) {
	frame := encodeFrame(opcode, payload)
	conn.Write(frame)
}

func (s *Server) sendClose(conn net.Conn, code uint16, reason string) {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, code)
	if reason != "" {
		payload = append(payload, []byte(reason)...)
	}
	s.writeFrame(conn, 0x08, payload)
}

func encodeFrame(opcode byte, payload []byte) []byte {
	fin := byte(0x80)
	length := len(payload)

	var header []byte
	if length < 126 {
		header = []byte{fin | opcode, byte(length)}
	} else if length < 65536 {
		header = make([]byte, 4)
		header[0] = fin | opcode
		header[1] = 126
		binary.BigEndian.PutUint16(header[2:4], uint16(length))
	} else {
		header = make([]byte, 10)
		header[0] = fin | opcode
		header[1] = 127
		binary.BigEndian.PutUint64(header[2:10], uint64(length))
	}

	result := make([]byte, len(header)+len(payload))
	copy(result, header)
	copy(result[len(header):], payload)
	return result
}

func (s *Server) handleMessage(text string, _ net.Conn) {
	var event map[string]any
	if err := json.Unmarshal([]byte(text), &event); err != nil {
		return
	}
	s.touchActivity()
	log.Printf("%s", text)

	if choice, ok := event["choice"]; ok && choice != nil {
		eventsPath := filepath.Join(s.cfg.StateDir, "events")
		f, err := os.OpenFile(eventsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(text + "\n")
			f.Close()
		}
	}
}

func (s *Server) broadcast(msg map[string]any) {
	data, _ := json.Marshal(msg)
	frame := encodeFrame(0x01, data)

	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for conn := range s.clients {
		_, err := conn.Write(frame)
		if err != nil {
			delete(s.clients, conn)
			conn.Close()
		}
	}
}

// ========== File Polling ==========

func (s *Server) pollFiles() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkFiles()
		case <-s.stopCh:
			return
		}
	}
}

func (s *Server) checkFiles() {
	entries, err := os.ReadDir(s.cfg.ContentDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".html") {
			continue
		}

		s.knownMu.Lock()
		wasKnown := s.knownFiles[e.Name()]
		if !wasKnown {
			s.knownFiles[e.Name()] = true
		}
		s.knownMu.Unlock()

		if !wasKnown {
			s.touchActivity()
			// Clear events file on new screen
			eventsPath := filepath.Join(s.cfg.StateDir, "events")
			os.Remove(eventsPath)
			log.Printf(`{"type":"screen-added","file":%q}`, filepath.Join(s.cfg.ContentDir, e.Name()))
			s.broadcast(map[string]any{"type": "reload"})
		}
	}
}

// ========== Helpers ==========

func (s *Server) newestScreen() string {
	entries, err := os.ReadDir(s.cfg.ContentDir)
	if err != nil {
		return ""
	}

	var newest string
	var newestTime time.Time
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".html") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newest = filepath.Join(s.cfg.ContentDir, e.Name())
		}
	}
	return newest
}

func isFullDocument(html string) bool {
	trimmed := strings.TrimLeft(html, " \t\r\n")
	lower := strings.ToLower(trimmed)
	return strings.HasPrefix(lower, "<!doctype") || strings.HasPrefix(lower, "<html")
}

func wrapInFrame(frameHTML, content string) string {
	return strings.Replace(frameHTML, "<!-- CONTENT -->", content, 1)
}

// ReadFileToString reads a file and returns its content as string.
func ReadFileToString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// CopyFile copies a file from src to dst.
func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
