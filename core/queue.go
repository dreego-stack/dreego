package core

import (
	"context"
	"time"
)

// JobHandler processes a single job. Implementations must respect ctx
// cancellation and return an error when the job could not be processed.
type JobHandler func(ctx context.Context, job Job) error

// JobMiddleware wraps a JobHandler. Middlewares are applied FIFO: the first
// registered middleware is the outermost wrapper.
type JobMiddleware func(next JobHandler) JobHandler

// Job is an opaque unit of work. ID is unique per caller, Name routes the
// job to the worker registered for it, Payload carries opaque bytes.
type Job struct {
	ID      string
	Name    string
	Payload []byte
}

// Queue is a background job queue contract, like database/sql: core defines
// the interface, plugins implement it (Redis, NATS, in-memory, ...). Core
// code stays transport-agnostic.
type Queue interface {
	// Dispatch enqueues job for immediate execution by the worker registered
	// for job.Name. It returns an error if no worker is registered or the
	// enqueue fails. It respects ctx cancellation.
	Dispatch(ctx context.Context, job Job) error
	// DispatchAfter enqueues job for execution after delay (delayed dispatch).
	// It respects ctx cancellation.
	DispatchAfter(ctx context.Context, job Job, delay time.Duration) error
	// DispatchBatch enqueues all jobs atomically (all-or-nothing, batching).
	// It respects ctx cancellation.
	DispatchBatch(ctx context.Context, jobs []Job) error
	// Worker registers handler for a job name. Registering a name twice is an
	// error. Middlewares registered via Use before this call wrap the handler
	// FIFO (first registered = outermost).
	Worker(name string, handler JobHandler) error
	// Use appends job middlewares. They wrap handlers FIFO (first registered
	// = outermost) and apply to all workers registered after Use.
	Use(middlewares ...JobMiddleware)
}
