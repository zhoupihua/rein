package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/visual"
)

var (
	visualHost    string
	visualURLHost string
	visualPort    int
)

var visualCmd = &cobra.Command{
	Use:   "visual",
	Short: "Visual brainstorming companion server",
}

var visualStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the visual brainstorming server",
	Run:   runVisualStart,
}

var visualStopCmd = &cobra.Command{
	Use:   "stop [state-dir]",
	Short: "Stop a running visual brainstorming server",
	Args:  cobra.MaximumNArgs(1),
	Run:   runVisualStop,
}

func init() {
	visualStartCmd.Flags().StringVar(&visualHost, "host", "127.0.0.1", "bind address")
	visualStartCmd.Flags().StringVar(&visualURLHost, "url-host", "", "hostname in returned URL (defaults to localhost when host is 127.0.0.1)")
	visualStartCmd.Flags().IntVar(&visualPort, "port", 0, "port (0 = auto-assign)")

	visualCmd.AddCommand(visualStartCmd)
	visualCmd.AddCommand(visualStopCmd)
	rootCmd.AddCommand(visualCmd)
}

func runVisualStart(cmd *cobra.Command, args []string) {
	projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	if projectDir == "" {
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			exitError("cannot determine project directory: %v", err)
		}
	}

	sessionDir := filepath.Join(projectDir, ".rein", "brainstorm")
	contentDir := filepath.Join(sessionDir, "content")
	stateDir := filepath.Join(sessionDir, "state")

	// Load template files from the skills directory
	framePath := findSkillFile("skills/define/frame-template.html")
	helperPath := findSkillFile("skills/define/helper.js")

	frameHTML := visual.ReadFileToString(framePath)
	if frameHTML == "" {
		frameHTML = defaultFrameTemplate()
	}
	helperJS := visual.ReadFileToString(helperPath)
	if helperJS == "" {
		helperJS = defaultHelperJS()
	}

	// Resolve owner PID
	ownerPID := os.Getppid()

	cfg := visual.Config{
		Port:       visualPort,
		Host:       visualHost,
		URLHost:    visualURLHost,
		ContentDir: contentDir,
		StateDir:   stateDir,
		OwnerPID:   ownerPID,
		FrameHTML:  frameHTML,
		HelperJS:   helperJS,
	}

	if isJSON() {
		// In JSON mode, start server in background and return info
		startInBackground(cfg)
		return
	}

	srv := visual.NewServer(cfg)
	if err := srv.Start(); err != nil {
		exitError("server error: %v", err)
	}
}

func startInBackground(cfg visual.Config) {
	exe, err := os.Executable()
	if err != nil {
		exitError("cannot find executable: %v", err)
	}

	args := []string{"visual", "start",
		"--host", cfg.Host,
		"--port", strconv.Itoa(cfg.Port),
	}
	if cfg.URLHost != "" {
		args = append(args, "--url-host", cfg.URLHost)
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = backgroundProcessAttr()

	if err := cmd.Start(); err != nil {
		exitError("failed to start server: %v", err)
	}

	// Wait briefly for server-info to appear
	infoPath := filepath.Join(cfg.StateDir, "server-info")
	for i := 0; i < 50; i++ {
		data, err := os.ReadFile(infoPath)
		if err == nil {
			fmt.Println(string(data))
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	exitError("server did not start within 5 seconds")
}

func runVisualStop(cmd *cobra.Command, args []string) {
	projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	if projectDir == "" {
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			exitError("cannot determine project directory: %v", err)
		}
	}

	stateDir := filepath.Join(projectDir, ".rein", "brainstorm", "state")
	if len(args) > 0 {
		stateDir = args[0]
	}

	infoPath := filepath.Join(stateDir, "server-info")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		exitError("no server-info found at %s", infoPath)
	}

	var info map[string]any
	if err := json.Unmarshal(data, &info); err != nil {
		exitError("invalid server-info: %v", err)
	}

	// Try to find and kill the process
	// The server-info doesn't include PID directly, find by port
	port, _ := info["port"].(float64)
	if port == 0 {
		exitError("no port in server-info")
	}

	if isJSON() {
		out, _ := json.Marshal(map[string]any{"type": "server-stop-requested", "port": int(port)})
		fmt.Println(string(out))
		return
	}

	fmt.Printf("Requesting server stop on port %d\n", int(port))
	fmt.Println("Send Ctrl+C to the running server process, or kill it manually.")
}

func findSkillFile(relPath string) string {
	// Check multiple possible locations for skill files
	candidates := []string{
		relPath,
		filepath.Join(".claude", relPath),
	}

	// Also check relative to the executable
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "..", "..", relPath),
		)
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return ""
}

func defaultFrameTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>rein Brainstorm</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    html, body { height: 100%; overflow: hidden; }
    :root {
      --bg-primary: #f5f5f7; --bg-secondary: #ffffff; --bg-tertiary: #e5e5e7;
      --border: #d1d1d6; --text-primary: #1d1d1f; --text-secondary: #86868b;
      --accent: #0071e3; --selected-bg: #e8f4fd; --selected-border: #0071e3;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg-primary: #1d1d1f; --bg-secondary: #2d2d2f; --bg-tertiary: #3d3d3f;
        --border: #424245; --text-primary: #f5f5f7; --text-secondary: #86868b;
        --accent: #0a84ff; --selected-bg: rgba(10,132,255,0.15); --selected-border: #0a84ff;
      }
    }
    body { font-family: system-ui, sans-serif; background: var(--bg-primary); color: var(--text-primary); display: flex; flex-direction: column; line-height: 1.5; }
    .header { background: var(--bg-secondary); padding: 0.5rem 1.5rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border); flex-shrink: 0; }
    .header h1 { font-size: 0.85rem; font-weight: 500; color: var(--text-secondary); }
    .header .status { font-size: 0.7rem; color: #34c759; display: flex; align-items: center; gap: 0.4rem; }
    .header .status::before { content: ''; width: 6px; height: 6px; background: #34c759; border-radius: 50%; }
    .main { flex: 1; overflow-y: auto; }
    #claude-content { padding: 2rem; min-height: 100%; }
    .indicator-bar { background: var(--bg-secondary); border-top: 1px solid var(--border); padding: 0.5rem 1.5rem; flex-shrink: 0; text-align: center; }
    .indicator-bar span { font-size: 0.75rem; color: var(--text-secondary); }
    .indicator-bar .selected-text { color: var(--accent); font-weight: 500; }
    .options { display: flex; flex-direction: column; gap: 0.75rem; }
    .option { background: var(--bg-secondary); border: 2px solid var(--border); border-radius: 12px; padding: 1rem 1.25rem; cursor: pointer; transition: all 0.15s ease; display: flex; align-items: flex-start; gap: 1rem; }
    .option:hover { border-color: var(--accent); }
    .option.selected { background: var(--selected-bg); border-color: var(--selected-border); }
    .option .letter { background: var(--bg-tertiary); color: var(--text-secondary); width: 1.75rem; height: 1.75rem; border-radius: 6px; display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 0.85rem; flex-shrink: 0; }
    .option.selected .letter { background: var(--accent); color: white; }
    .option .content { flex: 1; }
    .option .content h3 { font-size: 0.95rem; margin-bottom: 0.15rem; }
    .option .content p { color: var(--text-secondary); font-size: 0.85rem; margin: 0; }
  </style>
</head>
<body>
  <div class="header">
    <h1>rein Brainstorm</h1>
    <div class="status">Connected</div>
  </div>
  <div class="main">
    <div id="claude-content">
      <!-- CONTENT -->
    </div>
  </div>
  <div class="indicator-bar">
    <span id="indicator-text">Click an option above, then return to the terminal</span>
  </div>
</body>
</html>`
}

func defaultHelperJS() string {
	return `(function() {
  const WS_URL = 'ws://' + window.location.host;
  let ws = null;
  let eventQueue = [];
  function connect() {
    ws = new WebSocket(WS_URL);
    ws.onopen = () => { eventQueue.forEach(e => ws.send(JSON.stringify(e))); eventQueue = []; };
    ws.onmessage = (msg) => { const data = JSON.parse(msg.data); if (data.type === 'reload') window.location.reload(); };
    ws.onclose = () => { setTimeout(connect, 1000); };
  }
  function sendEvent(event) {
    event.timestamp = Date.now();
    if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(event));
    else eventQueue.push(event);
  }
  document.addEventListener('click', (e) => {
    const target = e.target.closest('[data-choice]');
    if (!target) return;
    sendEvent({ type: 'click', text: target.textContent.trim(), choice: target.dataset.choice, id: target.id || null });
    setTimeout(() => {
      const indicator = document.getElementById('indicator-text');
      if (!indicator) return;
      const container = target.closest('.options') || target.closest('.cards');
      const selected = container ? container.querySelectorAll('.selected') : [];
      if (selected.length === 0) indicator.textContent = 'Click an option above, then return to the terminal';
      else if (selected.length === 1) {
        const label = selected[0].querySelector('h3, .content h3, .card-body h3')?.textContent?.trim() || selected[0].dataset.choice;
        indicator.innerHTML = '<span class="selected-text">' + label + ' selected</span> — return to terminal to continue';
      } else indicator.innerHTML = '<span class="selected-text">' + selected.length + ' selected</span> — return to terminal to continue';
    }, 0);
  });
  window.selectedChoice = null;
  window.toggleSelect = function(el) {
    const container = el.closest('.options') || el.closest('.cards');
    const multi = container && container.dataset.multiselect !== undefined;
    if (container && !multi) container.querySelectorAll('.option, .card').forEach(o => o.classList.remove('selected'));
    if (multi) el.classList.toggle('selected'); else el.classList.add('selected');
    window.selectedChoice = el.dataset.choice;
  };
  window.brainstorm = { send: sendEvent, choice: (value, metadata = {}) => sendEvent({ type: 'choice', value, ...metadata }) };
  connect();
})();`
}

// backgroundProcessAttr returns platform-specific process attributes for background execution.
func backgroundProcessAttr() *syscall.SysProcAttr {
	if runtime.GOOS == "windows" {
		return &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x00000008, // DETACHED_PROCESS
		}
	}
	return &syscall.SysProcAttr{}
}


