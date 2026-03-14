package writer

import (
	"context"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Luvion1/mire/core"
	"github.com/Luvion1/mire/errors"
	"github.com/Luvion1/mire/util"
)

// LogProcessor defines the interface for the underlying logger that the AsyncLogger will use.
type LogProcessor interface {
	Log(ctx context.Context, level core.Level, msg []byte, keyvals ...[]byte)
	ErrorHandler() func(error)
	ErrOut() io.Writer
	ErrOutMu() *sync.Mutex
}

// logJob represents a logging job
type logJob struct {
	level   core.Level
	msg     []byte
	fields  map[string][]byte
	keyvals [][]byte
	ctx     context.Context
}

var logJobPool = sync.Pool{
	New: func() interface{} {
		return &logJob{}
	},
}

func getLogJob() *logJob {
	return logJobPool.Get().(*logJob)
}

func putLogJob(job *logJob) {
	job.msg = nil
	job.fields = nil
	job.keyvals = nil
	job.ctx = nil
	logJobPool.Put(job)
}

// AsyncLogger provides asynchronous logging to reduce latency
type AsyncLogger struct {
	processor                   LogProcessor
	logChan                     chan *logJob
	wg                          sync.WaitGroup
	workerCount                 int
	closed                      atomic.Bool
	logProcessTimeout           time.Duration
	disablePerLogContextTimeout bool
}

// NewAsyncLogger creates a new AsyncLogger
func NewAsyncLogger(processor LogProcessor, workerCount int, bufferSize int, logProcessTimeout time.Duration, disablePerLogContextTimeout bool) *AsyncLogger {
	al := &AsyncLogger{
		processor:                   processor,
		logChan:                     make(chan *logJob, bufferSize),
		workerCount:                 workerCount,
		logProcessTimeout:           logProcessTimeout,
		disablePerLogContextTimeout: disablePerLogContextTimeout,
	}

	for i := 0; i < workerCount; i++ {
		al.wg.Add(1)
		go al.worker()
	}

	return al
}

func (al *AsyncLogger) worker() {
	defer al.wg.Done()

	// Pre-allocate a batch of jobs to reduce channel operations if possible
	// But for now, we process one by one which is simpler

	for job := range al.logChan {
		al.processJob(job)
		putLogJob(job)
	}
}

func (al *AsyncLogger) processJob(job *logJob) {
	defer func() {
		if r := recover(); r != nil {
			if al.processor.ErrOut() != nil {
				mu := al.processor.ErrOutMu()
				mu.Lock()
				defer mu.Unlock()

				_, _ = al.processor.ErrOut().Write([]byte("recovering from panic in async logger worker: "))
				recoveredStr := util.ConvertValue(r)
				_, _ = al.processor.ErrOut().Write(util.StringToBytes(recoveredStr))
				_, _ = al.processor.ErrOut().Write([]byte("\n"))

				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				_, _ = al.processor.ErrOut().Write([]byte("stack trace: "))
				_, _ = al.processor.ErrOut().Write(buf[:n])
				_, _ = al.processor.ErrOut().Write([]byte("\n"))
			}
		}
	}()

	var ctx context.Context
	var cancel context.CancelFunc

	if al.logProcessTimeout > 0 && !al.disablePerLogContextTimeout {
		ctx, cancel = context.WithTimeout(job.ctx, al.logProcessTimeout)
	} else {
		ctx = job.ctx
	}

	// Call processor with appropriate arguments
	if job.keyvals != nil {
		al.processor.Log(ctx, job.level, job.msg, job.keyvals...)
	} else {
		// Convert fields to keyvals or call with empty keyvals
		if job.fields != nil && len(job.fields) > 0 {
			keyvals := make([][]byte, 0, len(job.fields)*2)
			for k, v := range job.fields {
				keyvals = append(keyvals, []byte(k), v)
			}
			al.processor.Log(ctx, job.level, job.msg, keyvals...)
		} else {
			// No fields or keyvals, just log the message
			al.processor.Log(ctx, job.level, job.msg)
		}
	}

	if cancel != nil {
		cancel()
	}
}

// LogZero queues a zero-allocation log job
func (al *AsyncLogger) LogZero(level core.Level, msg []byte, ctx context.Context, keyvals ...[]byte) {
	if al.closed.Load() {
		return
	}

	// Copy data if we're doing async
	msgCopy := make([]byte, len(msg))
	copy(msgCopy, msg)

	var keyvalsCopy [][]byte
	if len(keyvals) > 0 {
		keyvalsCopy = make([][]byte, len(keyvals))
		for i, kv := range keyvals {
			kvCopy := make([]byte, len(kv))
			copy(kvCopy, kv)
			keyvalsCopy[i] = kvCopy
		}
	}

	job := getLogJob()
	job.level = level
	job.msg = msgCopy
	job.keyvals = keyvalsCopy
	job.ctx = ctx

	select {
	case al.logChan <- job:
	default:
		putLogJob(job)
		if handler := al.processor.ErrorHandler(); handler != nil {
			handler(errors.ErrAsyncBufferFull)
		}
	}
}

// Log queues a log job for asynchronous processing
func (al *AsyncLogger) Log(level core.Level, msg []byte, fields map[string][]byte, ctx context.Context) {
	if al.closed.Load() {
		return
	}

	msgCopy := make([]byte, len(msg))
	copy(msgCopy, msg)

	// Fields are already copied or are safe?
	// Usually they are copied in the logger before calling asyncLogger.Log

	job := getLogJob()
	job.level = level
	job.msg = msgCopy
	job.fields = fields
	job.ctx = ctx

	select {
	case al.logChan <- job:
	default:
		putLogJob(job)
		if handler := al.processor.ErrorHandler(); handler != nil {
			handler(errors.ErrAsyncBufferFull)
		}
	}
}

// Close closes the async logger
func (al *AsyncLogger) Close() {
	if al.closed.CompareAndSwap(false, true) {
		close(al.logChan)
		al.wg.Wait()
	}
}
