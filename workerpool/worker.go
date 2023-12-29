package workerpool

type IWorker interface {
}

type Worker struct {
}

func NewAntsWorkerPool() IWorker {
	return &Worker{}
}
