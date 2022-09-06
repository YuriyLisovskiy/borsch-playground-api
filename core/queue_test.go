package core

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestQueue_Open_Success(t *testing.T) {
	q := NewQueue(1)
	err := q.Open(1)
	defer q.CloseAndWait(0 * time.Second)
	if err != nil {
		t.Errorf("Unable to open queue: %v", err)
	}
}

func TestQueue_Open_DoubleOpen(t *testing.T) {
	q := NewQueue(1)
	_ = q.Open(1)
	defer q.CloseAndWait(0 * time.Second)
	err := q.Open(1)
	if err == nil {
		t.Error("Queue was opened by mistake")
	}
}

type jobMock struct {
}

func (j *jobMock) Run() {
}

func TestQueue_Enqueue_Success(t *testing.T) {
	q := NewQueue(1)
	_ = q.Open(1)
	defer q.CloseAndWait(0 * time.Second)
	ok, err := q.Enqueue(&jobMock{})
	if err != nil {
		t.Errorf("Queue is not opened: %v", err)
	}

	if !ok {
		t.Error("Unable to enqueue job")
	}
}

func TestQueue_Enqueue_QueueIsNotOpened(t *testing.T) {
	q := NewQueue(1)
	_, err := q.Enqueue(&jobMock{})
	if err == nil {
		t.Error("Job was enqueued by mistake")
	}
}

func TestQueue_Enqueue_QueueIsFull(t *testing.T) {
	q := NewQueue(0)
	_ = q.Open(1)
	defer q.CloseAndWait(0 * time.Second)
	ok, _ := q.Enqueue(&jobMock{})
	if ok {
		t.Error("Job was enqueued by mistake")
	}
}

type benchJob struct {
	sleepTimeInSec int
}

func (j *benchJob) Run() {
	time.Sleep(time.Duration(j.sleepTimeInSec) * time.Second)
}

func TestQueue_Timings(t *testing.T) {
	expectedTimeInSec := 1
	workers := 5
	capacity := workers
	q := NewQueue(capacity)
	_ = q.Open(workers)

	start := time.Now()
	for i := 0; i < workers; i++ {
		ok, err := q.Enqueue(&benchJob{sleepTimeInSec: expectedTimeInSec})
		if err != nil {
			fmt.Println(fmt.Sprintf("err: %v, %d", err, i))
		}

		if !ok {
			fmt.Printf("not ok, %d\n", i)
		}
	}

	q.CloseAndWait(5 * time.Second)

	elapsed := time.Since(start)
	accuracy := 0.1
	if math.Abs(elapsed.Seconds()-float64(expectedTimeInSec)) > accuracy {
		t.Errorf(
			"Execution time is expected to be %ds with %f accuracy, actual %v",
			expectedTimeInSec,
			accuracy,
			elapsed,
		)
	}
}
