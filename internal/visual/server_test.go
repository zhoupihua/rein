package visual

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsFullDocument(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"<!DOCTYPE html><html></html>", true},
		{"  <!DOCTYPE html>", true},
		{"<html><body></body></html>", true},
		{"  <html>", true},
		{"<div>content</div>", false},
		{"Some text", false},
	}

	for _, tt := range tests {
		got := isFullDocument(tt.input)
		if got != tt.want {
			t.Errorf("isFullDocument(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestWrapInFrame(t *testing.T) {
	frame := `<html><body><!-- CONTENT --></body></html>`
	content := `<h2>Hello</h2>`
	result := wrapInFrame(frame, content)
	if result != `<html><body><h2>Hello</h2></body></html>` {
		t.Errorf("wrapInFrame produced unexpected result: %q", result)
	}
}

func TestNewestScreen(t *testing.T) {
	dir := t.TempDir()

	cfg := Config{
		Host:       "127.0.0.1",
		URLHost:    "localhost",
		ContentDir: dir,
		StateDir:   filepath.Join(dir, "state"),
	}
	srv := NewServer(cfg)

	// No files yet
	if got := srv.newestScreen(); got != "" {
		t.Errorf("expected empty, got %q", got)
	}

	// Create first file
	os.WriteFile(filepath.Join(dir, "a.html"), []byte("a"), 0644)
	time.Sleep(10 * time.Millisecond)

	// Create second file (newer)
	os.WriteFile(filepath.Join(dir, "b.html"), []byte("b"), 0644)

	got := srv.newestScreen()
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
	if filepath.Base(got) != "b.html" {
		t.Errorf("expected newest file b.html, got %q", filepath.Base(got))
	}
}

func TestHandleIndexWaitingPage(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		Host:       "127.0.0.1",
		URLHost:    "localhost",
		ContentDir: dir,
		StateDir:   filepath.Join(dir, "state"),
		FrameHTML:  "<html><body><!-- CONTENT --></body></html>",
		HelperJS:   "// helper",
	}
	srv := NewServer(cfg)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if body == "" {
		t.Error("expected non-empty response")
	}
}

func TestHandleIndexWithContent(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		Host:       "127.0.0.1",
		URLHost:    "localhost",
		ContentDir: dir,
		StateDir:   filepath.Join(dir, "state"),
		FrameHTML:  "<html><body><!-- CONTENT --></body></html>",
		HelperJS:   "// helper",
	}
	srv := NewServer(cfg)

	// Write a content fragment
	content := `<h2>Pick an option</h2>`
	os.WriteFile(filepath.Join(dir, "screen.html"), []byte(content), 0644)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.handleIndex(w, req)

	body := w.Body.String()
	if !contains(body, "<h2>Pick an option</h2>") {
		t.Error("response should contain the content")
	}
	if !contains(body, "// helper") {
		t.Error("response should contain the helper script")
	}
}

func TestHandleFileNotFound(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		Host:       "127.0.0.1",
		URLHost:    "localhost",
		ContentDir: dir,
		StateDir:   filepath.Join(dir, "state"),
	}
	srv := NewServer(cfg)

	req := httptest.NewRequest("GET", "/files/nonexistent.html", nil)
	w := httptest.NewRecorder()
	srv.handleFile(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleFileFound(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		Host:       "127.0.0.1",
		URLHost:    "localhost",
		ContentDir: dir,
		StateDir:   filepath.Join(dir, "state"),
	}
	srv := NewServer(cfg)

	os.WriteFile(filepath.Join(dir, "test.html"), []byte("<h1>Hi</h1>"), 0644)

	req := httptest.NewRequest("GET", "/files/test.html", nil)
	w := httptest.NewRecorder()
	srv.handleFile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Errorf("expected text/html content type, got %q", ct)
	}
}

func TestEncodeFrame(t *testing.T) {
	// Short payload
	frame := encodeFrame(0x01, []byte("hello"))
	if len(frame) < 2+5 {
		t.Errorf("frame too short: %d bytes", len(frame))
	}
	// Check FIN bit set and opcode
	if frame[0] != 0x81 {
		t.Errorf("expected FIN+TEXT opcode 0x81, got 0x%02x", frame[0])
	}
	if frame[1] != 5 {
		t.Errorf("expected length 5, got %d", frame[1])
	}
}

func TestDecodeFrame(t *testing.T) {
	// Build a masked client frame manually (decodeFrame requires masked frames)
	original := []byte("test message")
	mask := []byte{0x12, 0x34, 0x56, 0x78}
	masked := make([]byte, len(original))
	for i, b := range original {
		masked[i] = b ^ mask[i%4]
	}

	// Build frame: FIN+TEXT, mask bit set, length, mask key, masked payload
	frame := []byte{0x81, 0x80 | byte(len(original))}
	frame = append(frame, mask...)
	frame = append(frame, masked...)

	result, consumed, err := decodeFrame(frame)
	if err != nil {
		t.Fatalf("decodeFrame failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if consumed != len(frame) {
		t.Errorf("expected consumed=%d, got %d", len(frame), consumed)
	}
	if result.opcode != 0x01 {
		t.Errorf("expected opcode 0x01, got 0x%02x", result.opcode)
	}
	if string(result.payload) != string(original) {
		t.Errorf("expected payload %q, got %q", original, result.payload)
	}
}

func TestDecodeFrameIncomplete(t *testing.T) {
	// Too short to decode
	_, consumed, err := decodeFrame([]byte{0x81})
	if err != nil || consumed != 0 {
		t.Error("expected no error and 0 consumed for incomplete frame")
	}
}

func TestFilePolling(t *testing.T) {
	dir := t.TempDir()
	stateDir := filepath.Join(dir, "state")
	os.MkdirAll(stateDir, 0755)

	cfg := Config{
		Host:       "127.0.0.1",
		URLHost:    "localhost",
		ContentDir: dir,
		StateDir:   stateDir,
	}
	srv := NewServer(cfg)

	// No HTML files yet
	srv.checkFiles()
	if len(srv.knownFiles) != 0 {
		t.Error("expected no known files")
	}

	// Add an HTML file
	os.WriteFile(filepath.Join(dir, "screen1.html"), []byte("<h1>Test</h1>"), 0644)

	// Create events file so polling can clear it
	os.WriteFile(filepath.Join(stateDir, "events"), []byte("{}\n"), 0644)

	srv.checkFiles()
	if !srv.knownFiles["screen1.html"] {
		t.Error("expected screen1.html to be known")
	}

	// Events file should be cleared
	if _, err := os.Stat(filepath.Join(stateDir, "events")); !os.IsNotExist(err) {
		t.Error("expected events file to be removed on new screen")
	}
}

func TestStartInfoJSON(t *testing.T) {
	info := StartInfo{
		Type:      "server-started",
		Port:      52341,
		Host:      "127.0.0.1",
		URLHost:   "localhost",
		URL:       "http://localhost:52341",
		ScreenDir: "/tmp/content",
		StateDir:  "/tmp/state",
	}
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed StartInfo
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if parsed.Port != 52341 {
		t.Errorf("expected port 52341, got %d", parsed.Port)
	}
	if parsed.URL != "http://localhost:52341" {
		t.Errorf("unexpected URL: %q", parsed.URL)
	}
}

func TestMimeTypes(t *testing.T) {
	tests := map[string]string{
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript; charset=utf-8",
		".png":  "image/png",
		".svg":  "image/svg+xml",
		".xyz":  "application/octet-stream",
	}

	for ext, expected := range tests {
		got, ok := mimeTypes[ext]
		if ext == ".xyz" {
			if ok {
				t.Error("unexpected mime type for .xyz")
			}
			continue
		}
		if got != expected {
			t.Errorf("mimeTypes[%q] = %q, want %q", ext, got, expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
