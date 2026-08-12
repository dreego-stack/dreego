package core

import (
	"context"
	"errors"
	"testing"
	"time"
)

var _ Queue = (*fakeQueue)(nil)

type fakeQueue struct {
	workers     map[string]JobHandler
	middlewares []JobMiddleware
	afterCalls  []afterCall
}

type afterCall struct {
	job   Job
	delay time.Duration
}

func newFakeQueue() *fakeQueue {
	return &fakeQueue{workers: make(map[string]JobHandler)}
}

func (q *fakeQueue) Dispatch(ctx context.Context, job Job) error {
	h, ok := q.workers[job.Name]
	if !ok {
		return errors.New("queue: no worker for job name " + job.Name)
	}
	return h(ctx, job)
}

func (q *fakeQueue) DispatchAfter(ctx context.Context, job Job, delay time.Duration) error {
	q.afterCalls = append(q.afterCalls, afterCall{job: job, delay: delay})
	return nil
}

func (q *fakeQueue) DispatchBatch(ctx context.Context, jobs []Job) error {
	for _, job := range jobs {
		if err := q.Dispatch(ctx, job); err != nil {
			return err
		}
	}
	return nil
}

func (q *fakeQueue) Worker(name string, handler JobHandler) error {
	if _, ok := q.workers[name]; ok {
		return errors.New("queue: worker already registered for " + name)
	}
	q.workers[name] = wrap(handler, q.middlewares)
	return nil
}

func (q *fakeQueue) Use(middlewares ...JobMiddleware) {
	q.middlewares = append(q.middlewares, middlewares...)
}

func wrap(h JobHandler, middlewares []JobMiddleware) JobHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func TestQueueDispatchRunsRegisteredHandler(t *testing.T) {
	q := newFakeQueue()
	var got Job
	q.Worker("email", func(_ context.Context, job Job) error {
		got = job
		return nil
	})

	want := Job{ID: "j1", Name: "email", Payload: []byte("hello")}
	if err := q.Dispatch(context.Background(), want); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if got.ID != want.ID || got.Name != want.Name || string(got.Payload) != string(want.Payload) {
		t.Fatalf("handler received %+v, want %+v", got, want)
	}
}

func TestQueueDispatchNoWorkerReturnsError(t *testing.T) {
	q := newFakeQueue()
	if err := q.Dispatch(context.Background(), Job{Name: "unknown"}); err == nil {
		t.Fatal("dispatch without worker returned nil error, want error")
	}
}

func TestQueueWorkerDuplicateNameReturnsError(t *testing.T) {
	q := newFakeQueue()
	h := func(context.Context, Job) error { return nil }
	if err := q.Worker("email", h); err != nil {
		t.Fatalf("first worker registration: %v", err)
	}
	if err := q.Worker("email", h); err == nil {
		t.Fatal("duplicate worker registration returned nil error, want error")
	}
}

func TestQueueDispatchAfterPassesDelay(t *testing.T) {
	q := newFakeQueue()
	job := Job{ID: "j2", Name: "report", Payload: []byte("data")}
	want := 5 * time.Second

	if err := q.DispatchAfter(context.Background(), job, want); err != nil {
		t.Fatalf("dispatch after: %v", err)
	}
	if len(q.afterCalls) != 1 {
		t.Fatalf("recorded %d dispatch-after calls, want 1", len(q.afterCalls))
	}
	call := q.afterCalls[0]
	if call.delay != want {
		t.Fatalf("delay = %v, want %v", call.delay, want)
	}
	if call.job.ID != job.ID || call.job.Name != job.Name || string(call.job.Payload) != string(job.Payload) {
		t.Fatalf("recorded job %+v, want %+v", call.job, job)
	}
}

func TestQueueDispatchBatchDispatchesAll(t *testing.T) {
	q := newFakeQueue()
	var got []string
	q.Worker("job", func(_ context.Context, job Job) error {
		got = append(got, job.ID)
		return nil
	})

	jobs := []Job{{ID: "a", Name: "job"}, {ID: "b", Name: "job"}, {ID: "c", Name: "job"}}
	if err := q.DispatchBatch(context.Background(), jobs); err != nil {
		t.Fatalf("dispatch batch: %v", err)
	}
	if len(got) != len(jobs) {
		t.Fatalf("handler ran %d times, want %d", len(got), len(jobs))
	}
	for i, id := range got {
		if id != jobs[i].ID {
			t.Fatalf("dispatch order: got %q at %d, want %q", id, i, jobs[i].ID)
		}
	}
}

func TestQueueMiddlewareWrapsHandlerFIFO(t *testing.T) {
	q := newFakeQueue()
	var order []string

	q.Worker("plain", func(_ context.Context, job Job) error {
		order = append(order, "plain")
		return nil
	})

	q.Use(
		func(next JobHandler) JobHandler {
			return func(ctx context.Context, job Job) error {
				order = append(order, "m1:before")
				err := next(ctx, job)
				order = append(order, "m1:after")
				return err
			}
		},
		func(next JobHandler) JobHandler {
			return func(ctx context.Context, job Job) error {
				order = append(order, "m2:before")
				err := next(ctx, job)
				order = append(order, "m2:after")
				return err
			}
		},
	)
	q.Worker("wrapped", func(_ context.Context, job Job) error {
		order = append(order, "handler")
		return nil
	})

	if err := q.Dispatch(context.Background(), Job{Name: "plain"}); err != nil {
		t.Fatalf("dispatch plain: %v", err)
	}
	if err := q.Dispatch(context.Background(), Job{Name: "wrapped"}); err != nil {
		t.Fatalf("dispatch wrapped: %v", err)
	}

	want := []string{"plain", "m1:before", "m2:before", "handler", "m2:after", "m1:after"}
	if len(order) != len(want) {
		t.Fatalf("call order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("call order = %v, want %v", order, want)
		}
	}
}
