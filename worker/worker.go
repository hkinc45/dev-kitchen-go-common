package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Config holds the configuration for the pull subscriber worker pool.
type Config struct {
	StreamName          string
	Subject             string
	DurableName         string
	BatchSize           int
	MaxConcurrent       int
	MaxWait             time.Duration
	Handler             Handler
	JetStream           nats.JetStreamContext
}

// Handler is an interface that processing logic must implement.
type Handler interface {
	// Process handles a single NATS message.
	Process(ctx context.Context, msg *nats.Msg) error
	// GetLockingKey extracts a string key from a message to ensure sequential processing for the same resource.
	// If no specific locking is needed, it can return an empty string.
	GetLockingKey(msg *nats.Msg) (string, error)
}

// PullSubscriber manages a pool of workers to process messages from a NATS JetStream pull subscription.
type PullSubscriber struct {
	config    Config
	sub       *nats.Subscription
	mu        sync.Mutex
	active    bool
	keyLocks  map[string]*sync.Mutex
	keyLocksMu sync.RWMutex
	semaphore chan struct{}
}

// NewPullSubscriber creates and starts a new concurrent pull subscriber.
func NewPullSubscriber(cfg Config) (*PullSubscriber, error) {
	// Set sane defaults
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = 25
	}
	if cfg.MaxWait == 0 {
		cfg.MaxWait = 30 * time.Second
	}

	// Create the JetStream consumer
	_, err := cfg.JetStream.AddConsumer(cfg.StreamName, &nats.ConsumerConfig{
		Durable:       cfg.DurableName,
		AckPolicy:     nats.AckExplicitPolicy,
		FilterSubject: cfg.Subject,
		MaxDeliver:    5, // This is a reasonable default
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer for subject %s: %w", cfg.Subject, err)
	}

	// Create the pull subscription
	sub, err := cfg.JetStream.PullSubscribe(cfg.Subject, cfg.DurableName, nats.BindStream(cfg.StreamName))
	if err != nil {
		return nil, fmt.Errorf("failed to pull subscribe to subject %s: %w", cfg.Subject, err)
	}

	ps := &PullSubscriber{
		config:    cfg,
		sub:       sub,
		active:    true,
		keyLocks:  make(map[string]*sync.Mutex),
		semaphore: make(chan struct{}, cfg.MaxConcurrent),
	}

	go ps.startDispatcher()

	log.Printf("Successfully started concurrent subscriber for subject '%s' with durable name '%s'", cfg.Subject, cfg.DurableName)
	return ps, nil
}

// startDispatcher is the main loop that fetches messages and dispatches them to workers.
func (ps *PullSubscriber) startDispatcher() {
	for {
		ps.mu.Lock()
		if !ps.active {
			ps.mu.Unlock()
			return
		}
		ps.mu.Unlock()

		msgs, err := ps.sub.Fetch(ps.config.BatchSize, nats.MaxWait(ps.config.MaxWait))
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			log.Printf("ERROR: Failed to fetch messages for subject %s: %v", ps.config.Subject, err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, msg := range msgs {
			ps.semaphore <- struct{}{} // Acquire semaphore slot
			go ps.processMessage(msg)
		}
	}
}

// processMessage handles the full lifecycle of a single message, including locking and acknowledgement.
func (ps *PullSubscriber) processMessage(msg *nats.Msg) {
	defer func() {
		<-ps.semaphore // Release semaphore slot
	}()

	lockingKey, err := ps.config.Handler.GetLockingKey(msg)
	if err != nil {
		log.Printf("ERROR: Failed to get locking key for message on subject %s: %v. Naking message.", msg.Subject, err)
		_ = msg.NakWithDelay(5 * time.Second)
		return
	}

	// If a locking key is provided, acquire the specific lock for that key.
	if lockingKey != "" {
		keyMutex := ps.getKeyMutex(lockingKey)
		keyMutex.Lock()
		defer keyMutex.Unlock()
	}

	log.Printf("Processing message on subject %s with key '%s'", msg.Subject, lockingKey)
	
	// Create a context for the handler
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 5-minute timeout per message
	defer cancel()

	if err := ps.config.Handler.Process(ctx, msg); err != nil {
		log.Printf("ERROR: Handler failed to process message on subject %s: %v. Naking message.", msg.Subject, err)
		_ = msg.NakWithDelay(15 * time.Second) // Nak with a longer delay on processing failure
	} else {
		if err := msg.Ack(); err != nil {
			log.Printf("ERROR: Failed to ACK message on subject %s: %v", msg.Subject, err)
		} else {
			log.Printf("Successfully processed and ACKed message on subject %s with key '%s'", msg.Subject, lockingKey)
		}
	}
}

// getKeyMutex retrieves or creates a mutex for a specific key.
func (ps *PullSubscriber) getKeyMutex(key string) *sync.Mutex {
	ps.keyLocksMu.RLock()
	mutex, exists := ps.keyLocks[key]
	ps.keyLocksMu.RUnlock()

	if exists {
		return mutex
	}

	ps.keyLocksMu.Lock()
	// Double-check in case another goroutine created it while we were waiting for the write lock
	mutex, exists = ps.keyLocks[key]
	if !exists {
		mutex = &sync.Mutex{}
		ps.keyLocks[key] = mutex
	}
	ps.keyLocksMu.Unlock()
	
	return mutex
}

// Stop gracefully stops the subscriber.
func (ps *PullSubscriber) Stop() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if !ps.active {
		return
	}
	ps.active = false
	// Unsubscribe to stop receiving new messages
	if err := ps.sub.Unsubscribe(); err != nil {
		log.Printf("WARN: Error during unsubscribe for subject %s: %v", ps.config.Subject, err)
	}
	log.Printf("Stopped subscriber for subject '%s'", ps.config.Subject)
}
