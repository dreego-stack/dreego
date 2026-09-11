package wails_v3

import (
	"context"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TimerSnapshot struct {
	RemainingSeconds int    `json:"remainingSeconds"`
	Running          bool   `json:"running"`
	Status           string `json:"status"`
}

type TimerService struct {
	mu        sync.Mutex
	context   context.Context
	duration  int
	remaining int
	running   bool
	updatedAt time.Time
}

func NewTimerService(duration int) *TimerService {
	return &TimerService{duration: duration, remaining: duration}
}

func (t *TimerService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	t.mu.Lock()
	t.context = ctx
	t.updatedAt = time.Now()
	t.mu.Unlock()
	return nil
}

func (t *TimerService) ServiceShutdown() error {
	t.mu.Lock()
	t.update(time.Now())
	t.running = false
	t.context = nil
	t.mu.Unlock()
	return nil
}

func (t *TimerService) Snapshot() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.update(time.Now())
	return t.snapshot("Timer running.")
}

func (t *TimerService) Toggle() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	t.update(now)
	if t.running {
		t.running = false
		return t.snapshot("Timer paused.")
	}
	if t.remaining == 0 {
		t.remaining = t.duration
	}
	t.running = true
	t.updatedAt = now
	return t.snapshot("Timer running.")
}

func (t *TimerService) Reset() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.remaining = t.duration
	t.running = false
	t.updatedAt = time.Now()
	return t.snapshot("Timer reset to five minutes.")
}

func (t *TimerService) update(now time.Time) {
	if !t.running {
		return
	}
	elapsed := int(now.Sub(t.updatedAt) / time.Second)
	if elapsed == 0 {
		return
	}
	t.remaining -= elapsed
	t.updatedAt = t.updatedAt.Add(time.Duration(elapsed) * time.Second)
	if t.remaining <= 0 {
		t.remaining = 0
		t.running = false
	}
}

func (t *TimerService) snapshot(status string) TimerSnapshot {
	if t.remaining == 0 {
		status = "Focus session complete."
	} else if !t.running && status == "Timer running." {
		status = "Ready for a five-minute focus session."
	}
	return TimerSnapshot{RemainingSeconds: t.remaining, Running: t.running, Status: status}
}
