package workpool

import "sync"

func NewChannelPool(size int) *ChannelPool {
	return &ChannelPool{
		size:      size,
		queue:     make(chan Job, 100),
		waitGroup: &sync.WaitGroup{},
	}
}

type ChannelPool struct {
	size      int
	queue     chan Job
	waitGroup *sync.WaitGroup
	OnEnqueue func(Job)
	OnDequeue func(Job)
}

func (p ChannelPool) Start() {
	for i := 0; i < p.size; i++ {
		w := &ChannelWorker{id: i, pool: p}
		go w.Work()
	}
}

func (p ChannelPool) Stop() {
	close(p.queue)
	p.waitGroup.Wait()
}

func (p ChannelPool) Enqueue(j Job) error {
	p.queue <- j
	p.waitGroup.Add(1)
	if p.OnEnqueue != nil {
		p.OnEnqueue(j)
	}
	return nil
}

func (p ChannelPool) Dequeue(j Job) error {
	p.waitGroup.Done()
	if p.OnDequeue != nil {
		p.OnDequeue(j)
	}
	return nil
}

type ChannelWorker struct {
	id   int
	pool Pool
}

func (w ChannelWorker) Work() {
	for j := range w.pool.(ChannelPool).queue {
		j.Perform()
		w.pool.Dequeue(j)
	}
}
