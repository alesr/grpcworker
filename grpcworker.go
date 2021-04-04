package grpcworker

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

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
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func New(logger *zap.Logger, listener net.Listener, serviceDescription *grpc.ServiceDesc, service interface{}, opts ...Option) *Worker {
	w := &Worker{
		logger:   logger,
		listener: listener,
	}

	for _, opt := range opts {
		opt(w)
	}

	w.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(w.unaryInterceptors...),
		grpc.ChainStreamInterceptor(w.streamInterceptors...),
	)

	w.server.RegisterService(serviceDescription, service)
	return w
}

func (w *Worker) Init(*zap.Logger) error {
	w.logger.Named("grpc_server")

	return nil
}

func (w *Worker) Terminate() error {
	w.logger.Info("terminating grpc server")
	w.server.GracefulStop()
	return nil
}

func (w *Worker) Run() error {
	w.logger.Info("staring grpc server", zap.String("address", w.listener.Addr().String()))

	if err := w.server.Serve(w.listener); err != nil {
		return fmt.Errorf("could not serve grpc server: %s", err)
	}
	return nil
}
