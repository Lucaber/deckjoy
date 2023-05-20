package util

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

func ExecAsRoot(ctx context.Context, args ...string) error {
	polkitAgent := exec.CommandContext(ctx, "/usr/lib/polkit-kde-authentication-agent-1")
	polkitAgent.Stdout = NewLogWriter(log.DebugLevel)
	polkitAgent.Stderr = NewLogWriter(log.ErrorLevel)
	_ = polkitAgent.Start()
	defer func() {
		_ = polkitAgent.Process.Kill()
		_ = polkitAgent.Wait()
	}()

	time.Sleep(100 * time.Millisecond)

	err := Exec(ctx, "pkexec", args...)
	if err != nil {
		return err
	}

	return nil
}

func Exec(ctx context.Context, cmd string, args ...string) error {
	command := exec.CommandContext(ctx, cmd, args...)
	command.Cancel = func() error {
		err := command.Process.Kill()
		log.Infof("failed to kill process %v", err)
		return nil
	}
	command.Stdout = NewLogWriter(log.InfoLevel)
	command.Stderr = NewLogWriter(log.ErrorLevel)
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

type LogWriter struct {
	logger *log.Logger
	level  log.Level
}

func NewLogWriter(level log.Level) *LogWriter {
	lw := &LogWriter{
		level: level,
	}
	return lw
}

func (lw LogWriter) Write(b []byte) (n int, err error) {
	log.StandardLogger().Log(lw.level, fmt.Sprintf("process output: %s", string(b)))
	return len(b), nil
}
