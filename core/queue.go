/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package core

import (
	"errors"
	"sync"
	"time"
)

type QueueError string

func (e QueueError) Error() string {
	return string(e)
}

var (
	ErrQueueIsFull        QueueError = "queue is full"
	ErrQueueIsUnavailable QueueError = "enqueue is unavailable"
)

type Job interface {
	Run()
}

type Queue struct {
	jobChan  chan Job
	capacity int
	wg       *sync.WaitGroup
	isOpen   bool
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		capacity: capacity,
		wg:       &sync.WaitGroup{},
		isOpen:   false,
	}
}

func (q *Queue) Open(workersCount int) error {
	if q.isOpen {
		return errors.New("queue is already opened")
	}

	q.isOpen = true
	q.jobChan = make(chan Job, q.capacity)
	for i := 0; i < workersCount; i++ {
		q.wg.Add(1)
		go q.startWorker()
	}

	return nil
}

// Enqueue tries to enqueue a job to the given job channel. Returns true if
// the operation was successful, and false if enqueuing would have not been
// possible without blocking. Job is not enqueued in the latter case.
func (q *Queue) Enqueue(job Job) error {
	if !q.isOpen {
		return ErrQueueIsUnavailable
	}

	select {
	case q.jobChan <- job:
		return nil
	default:
		return ErrQueueIsFull
	}
}

func (q *Queue) Close() {
	if q.isOpen {
		close(q.jobChan)
	}
}

// Wait does a Wait on a sync.WaitGroup object but with a specified
// timeout. Returns true if the wait completed without timing out, false
// otherwise.
func (q *Queue) Wait(timeout time.Duration) bool {
	defer func() {
		q.isOpen = false
	}()

	ch := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (q *Queue) CloseAndWait(timeout time.Duration) {
	q.Close()
	q.Wait(timeout)
}

func (q *Queue) startWorker() {
	defer q.wg.Done()
	for job := range q.jobChan {
		job.Run()
	}
}
