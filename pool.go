package workpool

type Pool interface {
	Start()
	Stop()
	Enqueue(Job) error
	Dequeue(Job) error
}

type Job interface {
	Perform() error
}

type Worker interface {
	Work()
}
