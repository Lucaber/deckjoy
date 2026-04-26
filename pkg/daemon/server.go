package daemon

import (
	context "context"
	"github.com/lucaber/deckjoy/pkg/bluetooth"
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/lucaber/deckjoy/pkg/ipc"
	"github.com/lucaber/deckjoy/pkg/usb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"syscall"
	"time"
)

type Server struct {
	ipc.UnimplementedDeckJoyDaemonServer
	path      string
	server    *grpc.Server
	usb       *usb.USB
	bluetooth *bluetooth.Bluetooth
}

func (s *Server) InstallSudoers(ctx context.Context, empty *ipc.Empty) (*ipc.Empty, error) {
	filename := "/etc/sudoers.d/zzzdeckjoy"

	exe, err := os.Executable()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get executable path: %v", err)
	}

	err = os.WriteFile(filename, []byte("deck ALL=(ALL) NOPASSWD: "+exe+" daemon"), 0440)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "install sudoers failed: %v", err)
	}

	return &ipc.Empty{}, nil
}

func (s *Server) Stop(ctx context.Context, empty *ipc.Empty) (*ipc.Empty, error) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()
	return &ipc.Empty{}, nil
}

func (s *Server) InitUSB(ctx context.Context, request *ipc.Empty) (*ipc.Empty, error) {
	usb, err := usb.NewUSB()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup deck: %s", err.Error())
	}
	s.usb = usb

	_ = usb.Destroy()

	err = s.usb.SetupModules()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup kernel modules: %s", err.Error())
	}

	err = s.usb.SetupDeviceModules()
	if err != nil {
		log.WithError(err).Info("could not setup modules for usb device")
	}

	err = s.usb.SetupGadget()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup gadget: %s", err.Error())
	}

	return &ipc.Empty{}, nil
}

func (s *Server) SetupUSBJoystick(ctx context.Context, request *ipc.SetupJoystickRequest) (*ipc.SetupJoystickResponse, error) {
	if s.usb == nil {
		return nil, status.Error(codes.Unavailable, "gadget not setup")
	}
	path, err := s.usb.SetupJoystick(request.UserPermissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipc.SetupJoystickResponse{
		Path: path,
	}, nil
}

func (s *Server) SetupUSBKeyboard(ctx context.Context, request *ipc.SetupKeyboardRequest) (*ipc.SetupKeyboardResponse, error) {
	if s.usb == nil {
		return nil, status.Error(codes.Unavailable, "gadget not setup")
	}
	path, err := s.usb.SetupKeyboard(request.UserPermissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipc.SetupKeyboardResponse{
		Path: path,
	}, nil
}

func (s *Server) SetupUSBMouse(ctx context.Context, request *ipc.SetupMouseRequest) (*ipc.SetupMouseResponse, error) {
	if s.usb == nil {
		return nil, status.Error(codes.Unavailable, "gadget not setup")
	}
	path, err := s.usb.SetupMouse(request.UserPermissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipc.SetupMouseResponse{
		Path: path,
	}, nil
}

func (s *Server) InitBluetooth(ctx context.Context, request *ipc.Empty) (*ipc.Empty, error) {
	if err := bluetooth.EnsureBluetoothdInputDisabled(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not configure bluetoothd: %v", err)
	}

	s.bluetooth = bluetooth.NewBluetooth()

	descriptor := []byte{}
	descriptor = append(descriptor, hid.AddReportID(hid.JoystickReportDesc, 1)...)
	descriptor = append(descriptor, hid.AddReportID(hid.KeyboardReportDesc, 2)...)
	descriptor = append(descriptor, hid.AddReportID(hid.MouseReportDesc, 3)...)

	err := s.bluetooth.Run(ctx, bluetooth.SDPRecord{HIDDescriptor: descriptor})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup bluetooth: %v", err)
	}

	return &ipc.Empty{}, nil
}

func (s *Server) WriteBluetoothHIDData(ctx context.Context, request *ipc.WriteBluetoothHIDDataRequest) (*ipc.Empty, error) {
	if s.bluetooth == nil {
		return nil, status.Error(codes.Unavailable, "bluetooth not setup")
	}
	err := s.bluetooth.Write(request.Data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not send bluetooth: %v", err)
	}

	return &ipc.Empty{}, nil
}

var _ ipc.DeckJoyDaemonServer = &Server{}

func NewServer(path string) *Server {
	s := &Server{
		path:   path,
		server: grpc.NewServer(),
	}

	ipc.RegisterDeckJoyDaemonServer(s.server, s)

	return s
}

func (s *Server) Run() error {
	l, err := net.Listen("unix", s.path)
	defer func() {
		_ = os.Remove(s.path)
	}()
	if err != nil {
		return err
	}

	err = os.Chmod(s.path, os.ModePerm)
	if err != nil {
		return err
	}

	return s.server.Serve(l)
}

func (s *Server) Close() error {
	s.server.Stop()
	if s.usb != nil {
		err := s.usb.Destroy()
		if err != nil {
			return err
		}
		s.usb = nil
	}
	return nil
}
