package web

type Injector interface {
	Inject(fn interface{}) (interface{}, error)
}

type compilable interface {
	CompileDI(Injector)
}
