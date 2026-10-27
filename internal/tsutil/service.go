package tsutil

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"time"

	"tailscale.com/paths"
)

// IsTailscaledRunning checks if the tailscaled daemon is active or its socket is responding.
func IsTailscaledRunning(ctx context.Context) bool {
	// First check if the tailscale socket exists and is connectable.
	sock := paths.DefaultTailscaledSocket()
	if sock != "" {
		d := net.Dialer{Timeout: 500 * time.Millisecond}
		conn, err := d.DialContext(ctx, "unix", sock)
		if err == nil {
			_ = conn.Close()
			return true
		}
	}

	// Fallback to checking systemctl status
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", "tailscaled.service")
	err := cmd.Run()
	return err == nil
}

// EnsureTailscaledRunning starts the tailscaled systemd service using pkexec if not running.
func EnsureTailscaledRunning(ctx context.Context) error {
	if IsTailscaledRunning(ctx) {
		return nil
	}

	slog.Info("tailscaled is not active, attempting to start via pkexec systemctl...")

	// Launch pkexec to start the systemd unit
	cmd := exec.CommandContext(ctx, "pkexec", "systemctl", "start", "tailscaled.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("start tailscaled service: %w (output: %s)", err, string(out))
	}

	// Wait up to 5 seconds for the daemon and socket to become ready
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if IsTailscaledRunning(ctx) {
			slog.Info("tailscaled service successfully started and reachable")
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	return errors.New("tailscaled service started but socket is not reachable yet")
}

// StopTailscaledService stops the tailscaled systemd service using pkexec if it is running.
func StopTailscaledService(ctx context.Context) error {
	if !IsTailscaledRunning(ctx) {
		return nil
	}

	slog.Info("stopping tailscaled via pkexec systemctl...")
	cmd := exec.CommandContext(ctx, "pkexec", "systemctl", "stop", "tailscaled.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stop tailscaled service: %w (output: %s)", err, string(out))
	}
	return nil
}
