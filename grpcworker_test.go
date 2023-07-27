package grpcworker

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var _ net.Listener = (*fakeListener)(nil)

type fakeApp struct{}

type fakeListener struct{}

func (fl fakeListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (fl fakeListener) Close() error {
	return nil
}

func (fl fakeListener) Addr() net.Addr {
	return nil
}

func TestNew(t *testing.T) {
	logger := zap.NewNop()
	lis := fakeListener{}
	app := fakeApp{}
	serviceDesc := grpc.ServiceDesc{}

	t.Run("worker is instantiated with no options", func(t *testing.T) {
		t.Parallel()

		opts := []Option{}

		observedWorker := New(logger, lis, app, &serviceDesc, opts...)

		require.NotEmpty(t, observedWorker)

		assert.Equal(t, logger, observedWorker.logger)
		assert.Equal(t, lis, observedWorker.listener)
		assert.Equal(t, app, observedWorker.app)
		assert.Equal(t, &serviceDesc, observedWorker.serviceDesc)
		assert.Empty(t, observedWorker.unaryInterceptors)
		assert.Empty(t, observedWorker.streamInterceptors)
		assert.NotEmpty(t, observedWorker.server)
	})

	t.Run("worker is instantiated with unary interceptor option", func(t *testing.T) {
		t.Parallel()

		opts := []Option{
			WithUnaryInterceptor(nil),
		}

		observedWorker := New(logger, lis, app, &serviceDesc, opts...)

		require.NotEmpty(t, observedWorker)

		assert.NotEmpty(t, observedWorker.unaryInterceptors)
		assert.Empty(t, observedWorker.streamInterceptors)
	})

	t.Run("worker is instantiated with stream interceptor option", func(t *testing.T) {
		t.Parallel()

		opts := []Option{
			WithStreamInterceptor(nil),
		}

		observedWorker := New(logger, lis, app, &serviceDesc, opts...)

		require.NotEmpty(t, observedWorker)

		assert.Empty(t, observedWorker.unaryInterceptors)
		assert.NotEmpty(t, observedWorker.streamInterceptors)
	})

	t.Run("worker is instantiated with unary and stream interceptor options", func(t *testing.T) {
		t.Parallel()

		opts := []Option{
			WithUnaryInterceptor(nil),
			WithStreamInterceptor(nil),
		}

		observedWorker := New(logger, lis, app, &serviceDesc, opts...)

		require.NotEmpty(t, observedWorker)

		assert.NotEmpty(t, observedWorker.unaryInterceptors)
		assert.NotEmpty(t, observedWorker.streamInterceptors)
	})
}

func TestWorker_Init(t *testing.T) {
	logger := zap.NewNop()
	lis := fakeListener{}
	app := fakeApp{}
	serviceDesc := grpc.ServiceDesc{}

	observedWorker := New(logger, lis, app, &serviceDesc)

	t.Run("worker is initialized", func(t *testing.T) {
		t.Parallel()

		err := observedWorker.Init(logger)

		require.NoError(t, err)
	})
}
