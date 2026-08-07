package sse

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"sync"
	"time"

	dreego "codeberg.org/dreego/dreego/core"
)

//go:embed assets
var assetsFS embed.FS

type hub struct {
	mu   sync.Mutex
	subs map[chan string]struct{}
}

func newHub() *hub {
	return &hub{subs: make(map[chan string]struct{})}
}

func (h *hub) subscribe() chan string {
	ch := make(chan string, 16)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *hub) unsubscribe(ch chan string) {
	h.mu.Lock()
	delete(h.subs, ch)
	h.mu.Unlock()
}

func (h *hub) broadcast(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs {
		select {
		case ch <- msg:
		default:
		}
	}
}

type SSEPlugin struct {
	hub *hub
}

func New() *SSEPlugin {
	return &SSEPlugin{hub: newHub()}
}

func (p *SSEPlugin) Name() string {
	return "sse"
}

func (p *SSEPlugin) RegisterRoutes() {
	dreego.Register("GET", "/sse", p.handleSSE)
}

func (p *SSEPlugin) Middlewares() []func(http.Handler) http.Handler {
	return nil
}

func (p *SSEPlugin) Assets() fs.FS {
	sub, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		panic(err)
	}
	return sub
}

func (p *SSEPlugin) OnStart(ctx context.Context) error {
	return nil
}

func (p *SSEPlugin) OnShutdown(ctx context.Context) error {
	return nil
}

func (p *SSEPlugin) Broadcast(msg string) {
	p.hub.broadcast(msg)
}

func (p *SSEPlugin) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := p.hub.subscribe()
	defer p.hub.unsubscribe(ch)

	ctx := r.Context()
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	write := func(line string) {
		_, _ = w.Write([]byte(line))
		flusher.Flush()
	}

	write(": connected\n\n")
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			write("data: " + msg + "\n\n")
		case <-heartbeat.C:
			write(": ping\n\n")
		}
	}
}
