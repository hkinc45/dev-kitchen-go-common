package worker

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	// Silence slog during tests to prevent flooding the output during benchmarks
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Process(ctx context.Context, msg *nats.Msg) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockHandler) GetLockingKey(msg *nats.Msg) (string, error) {
	args := m.Called(msg)
	return args.String(0), args.Error(1)
}

func TestProcessMessage(t *testing.T) {
	ps := &PullSubscriber{
		config: Config{
			Handler: &MockHandler{},
		},
		semaphore: make(chan struct{}, 10),
		keyLocks:  make(map[string]*sync.Mutex),
	}

	handler := ps.config.Handler.(*MockHandler)
	msg := &nats.Msg{Subject: "test", Data: []byte("hello")}

	t.Run("Successful Processing", func(t *testing.T) {
		ps.semaphore <- struct{}{}
		handler.On("GetLockingKey", msg).Return("key1", nil).Once()
		handler.On("Process", mock.Anything, msg).Return(nil).Once()
		
		// Note: We can't easily test Ack() without a real connection, 
		// but we can test that the handler is called.
		ps.processMessage(msg)
		
		handler.AssertExpectations(t)
	})
}

func BenchmarkWorkerPool(b *testing.B) {
	ps := &PullSubscriber{
		config: Config{
			Handler: &MockHandler{},
		},
		semaphore: make(chan struct{}, 100),
		keyLocks:  make(map[string]*sync.Mutex),
	}
	handler := ps.config.Handler.(*MockHandler)
	msg := &nats.Msg{Subject: "test", Data: []byte("hello")}
	
	handler.On("GetLockingKey", mock.Anything).Return("", nil)
	handler.On("Process", mock.Anything, mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.semaphore <- struct{}{}
		ps.processMessage(msg)
	}
}
