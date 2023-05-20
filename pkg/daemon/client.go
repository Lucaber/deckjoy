package daemon

import (
	"context"
	"github.com/lucaber/deckjoy/pkg/ipc"
	"github.com/lucaber/deckjoy/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"time"
)

func RunSelfAsRoot(ctx context.Context, args ...string) error {
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}

	execArgs := []string{exePath}
	execArgs = append(execArgs, args...)
	return util.ExecAsRoot(ctx, execArgs...)
}

func RunDaemonProcess(ctx context.Context) chan error {
	errChan := make(chan error)
	go func() {
		err := RunSelfAsRoot(ctx, "daemon")
		// todo: stdout/stderr to errChan
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()
	return errChan
}

func StopDaemonProcess(ctx context.Context, path string) error {
	client, err := NewClient(ctx, path)
	if err != nil {
		return err
	}
	_, err = client.Stop(ctx, &ipc.Empty{})
	if err != nil {
		return err
	}
	return nil
}

func NewClient(ctx context.Context, path string) (ipc.DeckJoyDaemonClient, error) {
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return net.DialTimeout("unix", addr, time.Second)
	}
	c, err := grpc.DialContext(
		ctx,
		path,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	client := ipc.NewDeckJoyDaemonClient(c)
	return client, nil
}
