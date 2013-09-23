package workpool

import (
	"log"
	"testing"
	"time"
)

var workDelay time.Duration

func init() {
	workDelay = time.Millisecond
}

type MockJob struct {
	Id    string
	State string
}

func NewMockJob(id string) *MockJob {
	return &MockJob{id, "pending"}
}

func (m *MockJob) Perform() error {
	log.Print("working...")
	time.Sleep(workDelay) // simulate work
	m.State = "done"
	log.Print("done!")
	return nil
}

func TestCallbacks(t *testing.T) {
	p := NewChannelPool(2)
	defer p.Stop()
	p.OnEnqueue = func(j Job) {
		mj := j.(*MockJob)
		if mj.Id != "test" {
			t.Error("wrong job", mj)
		}
		if mj.State != "pending" {
			t.Errorf("expected pending, got %s", mj.State)
		}
	}
	p.OnDequeue = func(j Job) {
		mj := j.(*MockJob)
		if mj.Id != "test" {
			t.Error("wrong job", mj)
		}
		if mj.State != "done" {
			t.Errorf("expected done, got %s", mj.State)
		}
	}
	p.Enqueue(NewMockJob("test"))
	p.Start()
}

func TestSerial(t *testing.T) {
	p := NewChannelPool(1)
	p.Enqueue(NewMockJob("1"))
	p.Enqueue(NewMockJob("2"))
	start := time.Now()
	p.Start()
	p.Stop()

	if time.Since(start) < (workDelay * 2) {
		t.Error("job not worked serially")
	}
}

func TestParallel(t *testing.T) {
	p := NewChannelPool(2)
	p.Enqueue(NewMockJob("1"))
	p.Enqueue(NewMockJob("2"))
	start := time.Now()
	p.Start()
	p.Stop()

	window := (workDelay + workDelay/4) // 25% overhead
	if time.Since(start) > window {
		t.Error("job not worked in parallel", time.Since(start), window)
	}
}
