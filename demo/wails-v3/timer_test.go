package wails_v3

import (
	"context"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestTimerServiceLifecycleAndState(t *testing.T) {
	service := NewTimerService(300)
	ctx, cancel := context.WithCancel(context.Background())
	if err := service.ServiceStartup(ctx, application.ServiceOptions{}); err != nil {
		t.Fatalf("ServiceStartup: %v", err)
	}
	started := service.Toggle()
	if !started.Running || started.RemainingSeconds != 300 {
		t.Fatalf("started snapshot = %#v", started)
	}
	service.mu.Lock()
	service.updatedAt = service.updatedAt.Add(-2 * time.Second)
	service.mu.Unlock()
	if got := service.Snapshot().RemainingSeconds; got != 298 {
		t.Fatalf("remaining seconds = %d, want 298", got)
	}
	cancel()
	if err := service.ServiceShutdown(); err != nil {
		t.Fatalf("ServiceShutdown: %v", err)
	}
	service.mu.Lock()
	defer service.mu.Unlock()
	if service.context != nil || service.running {
		t.Fatalf("service retained lifecycle state: %#v", service)
	}
}
