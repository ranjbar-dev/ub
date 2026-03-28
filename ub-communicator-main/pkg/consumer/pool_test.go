package consumer

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"ub-communicator/pkg/messaging"
)

// mockMessagingService counts Send calls and optionally delays.
type mockMessagingService struct {
	sendCount int64 // atomic
	delay     time.Duration
}

func (m *mockMessagingService) Send(msg messaging.Message) error {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	atomic.AddInt64(&m.sendCount, 1)
	return nil
}

func (m *mockMessagingService) CreateMessage(data []byte) (messaging.Message, error) {
	return messaging.Message{Type: "EMAIL", Receiver: "test@test.com"}, nil
}

func TestPool_ProcessesAllMessages(t *testing.T) {
	ms := &mockMessagingService{}
	pool := NewPool(ms)
	collector := pool.StartDispatcher(3)

	const messageCount = 50
	var wg sync.WaitGroup
	wg.Add(messageCount)

	for i := 0; i < messageCount; i++ {
		go func() {
			defer wg.Done()
			work := Work{
				Message: messaging.Message{Type: "EMAIL", Receiver: "a@b.com"},
			}
			collector.Work <- work
		}()
	}

	wg.Wait()
	// Give workers time to process
	time.Sleep(100 * time.Millisecond)

	count := atomic.LoadInt64(&ms.sendCount)
	if count != messageCount {
		t.Errorf("processed %d messages, want %d", count, messageCount)
	}

	collector.End <- true
}

func TestPool_WorkDistribution(t *testing.T) {
	// With slow workers, work should distribute across multiple goroutines
	ms := &mockMessagingService{delay: 10 * time.Millisecond}
	pool := NewPool(ms)
	collector := pool.StartDispatcher(5)

	const messageCount = 20
	for i := 0; i < messageCount; i++ {
		collector.Work <- Work{
			Message: messaging.Message{Type: "EMAIL"},
		}
	}

	// Wait for all to process
	deadline := time.After(5 * time.Second)
	for {
		if atomic.LoadInt64(&ms.sendCount) >= messageCount {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("timeout: only %d of %d processed", atomic.LoadInt64(&ms.sendCount), messageCount)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	collector.End <- true
}

func TestPool_StopsOnEnd(t *testing.T) {
	ms := &mockMessagingService{}
	pool := NewPool(ms)
	collector := pool.StartDispatcher(3)

	// Send some work first
	for i := 0; i < 5; i++ {
		collector.Work <- Work{
			Message: messaging.Message{Type: "EMAIL"},
		}
	}
	time.Sleep(50 * time.Millisecond)

	// Signal end
	collector.End <- true

	// Workers should stop — verify no goroutine leak by attempting more work
	// After end, workers should not process new work
	time.Sleep(100 * time.Millisecond)
	before := atomic.LoadInt64(&ms.sendCount)

	// Try to send more work (this should not be processed since workers stopped)
	select {
	case collector.Work <- Work{Message: messaging.Message{Type: "EMAIL"}}:
		// Work was sent, but shouldn't be processed since workers stopped
	default:
		// Channel is full or closed, which is fine
	}

	time.Sleep(100 * time.Millisecond)
	after := atomic.LoadInt64(&ms.sendCount)

	if after != before {
		t.Errorf("workers still processing after end: before=%d, after=%d", before, after)
	}
}
