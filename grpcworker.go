package grpcworker

import (
	"fmt"
	"net"

	"github.com/voi-oss/svc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var _ svc.Worker = (*Worker)(nil)

type Option func(w *Worker)

func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) Option {
	return func(w *Worker) {
		w.unaryInterceptors = append(w.unaryInterceptors, interceptor)
	}
}

func WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) Option {
	return func(w *Worker) {
		w.streamInterceptors = append(w.streamInterceptors, interceptor)
	}
}

type Worker struct {
	logger             *zap.Logger
	listener           net.Listener
	server             *grpc.Server
	serviceDesc        *grpc.ServiceDesc
	app                any
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

// New creates a new instance of the Worker.
func New(logger *zap.Logger, lis net.Listener, app any, serviceDesc *grpc.ServiceDesc, opts ...Option) *Worker {
	w := &Worker{
		logger:      logger,
		listener:    lis,
		app:         app,
		serviceDesc: serviceDesc,
	}

	for _, opt := range opts {
		opt(w)
	}

	w.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(w.unaryInterceptors...),
		grpc.ChainStreamInterceptor(w.streamInterceptors...),
	)
	return w
}

// Init initializes the worker.
func (w *Worker) Init(logger *zap.Logger) error {
	w.logger.Named("grpc_worker")
	w.server.RegisterService(w.serviceDesc, w.app)
	return nil
}

// Run starts the gRPC server.
func (w *Worker) Run() error {
	w.logger.Info("starting grpc server", zap.String("address", w.listener.Addr().String()))
	if err := w.server.Serve(w.listener); err != nil {
		w.logger.Error("failed to serve grpc server", zap.Error(err))
		return fmt.Errorf("could not serve grpc server: %w", err)
	}
	return nil
}

// Terminate stops the gRPC server gracefully.
func (w *Worker) Terminate() error {
	w.logger.Info("terminating grpc server")
	w.server.GracefulStop()
	return nil
}
