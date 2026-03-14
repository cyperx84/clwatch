package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

func runService(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: clwatch service <install|uninstall|start|stop|status|logs>\n")
		return 1
	}

	sub := args[0]
	switch runtime.GOOS {
	case "darwin":
		return runServiceDarwin(sub)
	case "linux":
		return runServiceLinux(sub)
	default:
		fmt.Fprintf(os.Stderr, "service management not supported on %s\n", runtime.GOOS)
		fmt.Fprintf(os.Stderr, "Run clwatch watch manually: clwatch watch --interval 1h\n")
		return 1
	}
}

// --- macOS launchd ---

const launchdPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.clwatch.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.BinPath}}</string>
        <string>watch</string>
        <string>--interval</string>
        <string>1h</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.LogDir}}/clwatch.log</string>
    <key>StandardErrorPath</key>
    <string>{{.LogDir}}/clwatch-error.log</string>
    <key>ThrottleInterval</key>
    <integer>60</integer>
</dict>
</plist>
`

func launchdPlistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", "io.clwatch.agent.plist")
}

func launchdLogDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Logs", "clwatch")
}

func runServiceDarwin(sub string) int {
	plistPath := launchdPlistPath()
	logDir := launchdLogDir()

	switch sub {
	case "install":
		binPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot determine clwatch binary path: %v\n", err)
			return 1
		}
		binPath, _ = filepath.Abs(binPath)

		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating log dir: %v\n", err)
			return 1
		}

		f, err := os.Create(plistPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating plist: %v\n", err)
			return 1
		}
		defer f.Close()

		tmpl := template.Must(template.New("plist").Parse(launchdPlist))
		tmpl.Execute(f, map[string]string{"BinPath": binPath, "LogDir": logDir})

		if out, err := exec.Command("launchctl", "load", plistPath).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "error loading service: %v\n%s", err, out)
			return 1
		}

		fmt.Printf("✓ clwatch service installed and started\n")
		fmt.Printf("  Binary:  %s\n", binPath)
		fmt.Printf("  Plist:   %s\n", plistPath)
		fmt.Printf("  Logs:    %s/clwatch.log\n", logDir)
		fmt.Printf("  Interval: every 1 hour\n")

	case "uninstall":
		exec.Command("launchctl", "unload", plistPath).Run()
		if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error removing plist: %v\n", err)
			return 1
		}
		fmt.Printf("✓ clwatch service removed\n")

	case "start":
		if out, err := exec.Command("launchctl", "load", plistPath).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n%s", err, out)
			return 1
		}
		fmt.Printf("✓ clwatch service started\n")

	case "stop":
		if out, err := exec.Command("launchctl", "unload", plistPath).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n%s", err, out)
			return 1
		}
		fmt.Printf("✓ clwatch service stopped\n")

	case "status":
		out, err := exec.Command("launchctl", "list", "io.clwatch.agent").CombinedOutput()
		if err != nil {
			fmt.Printf("clwatch service: not running\n")
			return 0
		}
		if strings.Contains(string(out), "io.clwatch.agent") {
			fmt.Printf("clwatch service: running ✓\n")
			fmt.Printf("  Logs: %s/clwatch.log\n", logDir)
		} else {
			fmt.Printf("clwatch service: not found\n")
		}

	case "logs":
		logFile := filepath.Join(logDir, "clwatch.log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "no log file found at %s\n", logFile)
			return 1
		}
		out, _ := exec.Command("tail", "-50", logFile).CombinedOutput()
		fmt.Print(string(out))

	default:
		fmt.Fprintf(os.Stderr, "unknown service command: %s\n", sub)
		fmt.Fprintf(os.Stderr, "Usage: clwatch service <install|uninstall|start|stop|status|logs>\n")
		return 1
	}
	return 0
}

// --- Linux systemd ---

const systemdUnit = `[Unit]
Description=clwatch — AI coding tool update watcher
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart={{.BinPath}} watch --interval 1h
Restart=on-failure
RestartSec=30

[Install]
WantedBy=default.target
`

func systemdUnitPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "systemd", "user", "clwatch.service")
}

func runServiceLinux(sub string) int {
	unitPath := systemdUnitPath()

	switch sub {
	case "install":
		binPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot determine binary path: %v\n", err)
			return 1
		}
		binPath, _ = filepath.Abs(binPath)

		if err := os.MkdirAll(filepath.Dir(unitPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating systemd dir: %v\n", err)
			return 1
		}

		f, err := os.Create(unitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating unit file: %v\n", err)
			return 1
		}
		defer f.Close()

		tmpl := template.Must(template.New("unit").Parse(systemdUnit))
		tmpl.Execute(f, map[string]string{"BinPath": binPath})

		exec.Command("systemctl", "--user", "daemon-reload").Run()
		if out, err := exec.Command("systemctl", "--user", "enable", "--now", "clwatch").CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "error enabling service: %v\n%s", err, out)
			return 1
		}

		fmt.Printf("✓ clwatch service installed and started\n")
		fmt.Printf("  Binary: %s\n", binPath)
		fmt.Printf("  Unit:   %s\n", unitPath)
		fmt.Printf("  Logs:   journalctl --user -u clwatch -f\n")

	case "uninstall":
		exec.Command("systemctl", "--user", "disable", "--now", "clwatch").Run()
		os.Remove(unitPath)
		exec.Command("systemctl", "--user", "daemon-reload").Run()
		fmt.Printf("✓ clwatch service removed\n")

	case "start":
		if out, err := exec.Command("systemctl", "--user", "start", "clwatch").CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n%s", err, out)
			return 1
		}
		fmt.Printf("✓ clwatch service started\n")

	case "stop":
		if out, err := exec.Command("systemctl", "--user", "stop", "clwatch").CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n%s", err, out)
			return 1
		}
		fmt.Printf("✓ clwatch service stopped\n")

	case "status":
		out, _ := exec.Command("systemctl", "--user", "status", "clwatch").CombinedOutput()
		fmt.Print(string(out))

	case "logs":
		out, _ := exec.Command("journalctl", "--user", "-u", "clwatch", "-n", "50", "--no-pager").CombinedOutput()
		fmt.Print(string(out))

	default:
		fmt.Fprintf(os.Stderr, "unknown service command: %s\n", sub)
		return 1
	}
	return 0
}
