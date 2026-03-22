package worker

import (
	"context"
	"fmt"
	"goflow/logger"
	"goflow/retry"
	"goflow/service"
	"io"
	"log"
	"sync"
)

type Pool struct {
	numWorkers int
	taskChan   chan *ProcessingTask
	resultChan chan *ProcessingResult
	ErrChan    chan error
	wg         sync.WaitGroup
	consumer   service.EventConsumer
	downloader service.FileDownloader
	maxRetries int
	limiter    service.RateLimiter
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *logger.Logger
}

type ProcessingResult struct {
	Task  *ProcessingTask
	File  io.ReadCloser
	Error error
}

func NewPool(
	numWorkers int,
	consumer service.EventConsumer,
	downloader service.FileDownloader,
	maxRetries int,
	limiter service.RateLimiter,
) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	return &Pool{
		numWorkers: numWorkers,
		taskChan:   make(chan *ProcessingTask, numWorkers),
		resultChan: make(chan *ProcessingResult),
		ErrChan:    make(chan error),
		consumer:   consumer,
		downloader: downloader,
		maxRetries: maxRetries,
		limiter:    limiter,
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger.New(),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	p.wg.Add(1)
	go p.consumerLoop()
}

func (p *Pool) consumerLoop() {
	defer p.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("Consumer loop panicked: %v", r)
			log.Println(errMsg)
			p.ErrChan <- fmt.Errorf(errMsg)
		}
	}()

	for {
		select {
		case <-p.ctx.Done():
			close(p.taskChan)
			return
		default:
			event, err := p.consumer.Consume(p.ctx)
			if err != nil {
				p.ErrChan <- err
				continue
			}
			if event == nil {
				continue
			}

			task := NewProcessingTask(event)
			p.taskChan <- task
		}
	}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("Worker %d panicked: %v", id, r)
			log.Println(errMsg)
			p.ErrChan <- fmt.Errorf(errMsg)
		}
	}()

	for task := range p.taskChan {
		ctx := logger.WithCorrelationID(p.ctx, task.CorrelationID)

		p.logger.Info(ctx, "Processing started", task.Event.DocumentID)

		if err := p.limiter.Acquire(p.ctx); err != nil {
			p.ErrChan <- err
			continue
		}
		defer p.limiter.Release()

		var file io.ReadCloser
		var err error

		objectName := task.Event.ObjectName
		if objectName == "" {
			objectName = task.Event.Filename
		}

		err = retry.Retry(p.ctx, p.maxRetries, func() error {
			f, err := p.downloader.Download(p.ctx, task.Event.BucketName, objectName)
			file = f
			return err
		})

		result := &ProcessingResult{
			Task:  task,
			File:  file,
			Error: err,
		}
		p.resultChan <- result
	}
}

func (p *Pool) Results() <-chan *ProcessingResult {
	return p.resultChan
}

func (p *Pool) Errors() <-chan error {
	return p.ErrChan
}

func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()
	close(p.resultChan)
	close(p.ErrChan)
}
