package daemon

import (
	context "context"
	"github.com/lucaber/deckjoy/pkg/ipc"
	"github.com/lucaber/deckjoy/pkg/setup"
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
	path   string
	server *grpc.Server
	deck   *setup.Deck
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

func (s *Server) Init(ctx context.Context, request *ipc.Empty) (*ipc.Empty, error) {
	deck, err := setup.NewDeck()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup deck: %s", err.Error())
	}
	s.deck = deck

	_ = deck.Destroy()

	err = s.deck.SetupModules()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup kernel modules: %s", err.Error())
	}

	err = s.deck.SetupDeviceModules()
	if err != nil {
		log.WithError(err).Info("could not setup modules for usb device")
	}

	err = s.deck.SetupGadget()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not setup gadget: %s", err.Error())
	}

	return &ipc.Empty{}, nil
}

func (s *Server) SetupJoystick(ctx context.Context, request *ipc.SetupJoystickRequest) (*ipc.SetupJoystickResponse, error) {
	if s.deck == nil {
		return nil, status.Error(codes.Unavailable, "gadget not setup")
	}
	path, err := s.deck.SetupJoystick(request.UserPermissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipc.SetupJoystickResponse{
		Path: path,
	}, nil
}

func (s *Server) SetupKeyboard(ctx context.Context, request *ipc.SetupKeyboardRequest) (*ipc.SetupKeyboardResponse, error) {
	if s.deck == nil {
		return nil, status.Error(codes.Unavailable, "gadget not setup")
	}
	path, err := s.deck.SetupKeyboard(request.UserPermissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipc.SetupKeyboardResponse{
		Path: path,
	}, nil
}

func (s *Server) SetupMouse(ctx context.Context, request *ipc.SetupMouseRequest) (*ipc.SetupMouseResponse, error) {
	if s.deck == nil {
		return nil, status.Error(codes.Unavailable, "gadget not setup")
	}
	path, err := s.deck.SetupMouse(request.UserPermissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipc.SetupMouseResponse{
		Path: path,
	}, nil
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
	if s.deck != nil {
		err := s.deck.Destroy()
		if err != nil {
			return err
		}
		s.deck = nil
	}
	return nil
}
