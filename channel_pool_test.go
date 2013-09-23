package workpool

import "testing"

type MockJob struct {
	Id    string
	State string
}

func NewMockJob(id string) *MockJob {
	return &MockJob{id, "pending"}
}

func (m *MockJob) Perform() error {
	m.State = "done"
	return nil
}

func TestCallbacks(t *testing.T) {
	p := NewChannelPool(2)
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
	p.Stop()
}
