package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stage := stage
		in = stage(process(done, in))
	}
	return in
}

func process(done In, in In) Out {
	bi := make(Bi)
	go func() {
		defer close(bi)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case bi <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return bi
}
