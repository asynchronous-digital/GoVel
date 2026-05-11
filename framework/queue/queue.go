package queue
// Queue system (Phase 3)
// Background job processing with multiple drivers

package queue

import (
	"context"
	"time"
)

// Job represents a unit of work to be processed asynchronously
type Job interface {
	// Unique identifier for the job
	ID() string
	
	// Job name/type
	Name() string
	
	// Job payload
	Payload() map[string]interface{}
	
	// Handle the job
	Handle(ctx context.Context) error
	
	// Failed job handling
	Failed(err error) error
	
	// Job timeout
	Timeout() time.Duration
}

// Queue manages job dispatch and processing
type Queue interface {
	// Push a job onto the queue
	Push(ctx context.Context, job Job) error
	
	// Push a job with delay
	PushLater(ctx context.Context, job Job, delay time.Duration) error
	
	// Process jobs from the queue
	Process(ctx context.Context, maxJobs int) error
	
	// Get queue stats
	Count(ctx context.Context) (int64, error)
	
	// Clear the queue
	Clear(ctx context.Context) error
}

// Driver implementations
type Driver string

const (
	DriverSync     Driver = "sync"     // Process immediately (testing)
	DriverDatabase Driver = "database" // Store in database
	DriverRedis    Driver = "redis"    // Use Redis (future)
	DriverSQS      Driver = "sqs"      // AWS SQS (future)
)

// Worker processes jobs from the queue
type Worker interface {
	Handle(ctx context.Context, queue Queue) error
}

// Job examples:
/*
type SendWelcomeEmail struct {
    UserID uint
    Email  string
}

func (j SendWelcomeEmail) ID() string {
    return fmt.Sprintf("send-welcome-email-%d", j.UserID)
}

func (j SendWelcomeEmail) Name() string {
    return "send-welcome-email"
}

func (j SendWelcomeEmail) Handle(ctx context.Context) error {
    // Send email logic
    return nil
}

func (j SendWelcomeEmail) Timeout() time.Duration {
    return 30 * time.Second
}
*/

// Usage:
/*
// Push a job
queue.Push(ctx, &SendWelcomeEmail{
    UserID: user.ID,
    Email: user.Email,
})

// Push with delay
queue.PushLater(ctx, &SendWelcomeEmail{...}, 5 * time.Minute)

// Process queue (usually in background worker)
worker := queue.NewWorker()
worker.Handle(ctx, queue)
*/

// Processor handles job processing and retries
type Processor interface {
	Process(ctx context.Context, job Job) error
	Retry(ctx context.Context, job Job, attempt int) error
	Failed(ctx context.Context, job Job, err error) error
}

// Future implementation with retry logic, dead-letter queue, etc.
