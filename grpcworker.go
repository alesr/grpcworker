package grpcworker

import (
	"fmt"
	"net"

	"github.com/voi-oss/svc"
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

var _ svc.Worker = (*Worker)(nil)

type Worker struct {
	logger             *zap.Logger
	listener           net.Listener
	server             *grpc.Server
	serviceDesc        *grpc.ServiceDesc
	app                interface{}
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func New(logger *zap.Logger, lis net.Listener, app interface{}, serviceDesc *grpc.ServiceDesc, opts ...Option) *Worker {
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

func (w *Worker) Init(*zap.Logger) error {
	w.logger.Named("grpc_worker")
	w.server.RegisterService(w.serviceDesc, w.app)
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
