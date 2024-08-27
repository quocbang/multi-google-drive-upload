package workerpool

type IWorker interface {
}

type Ants struct {
}

func NewAntsWorkerPool() IWorker {
	return &Ants{}
}
