package rest

type handler func()

type middleware struct {
	task handler
}

type route struct {
	method  string
	pattern string
	task    handler
}

type exception struct {
	code string
	task handler
}

type list struct {
	middlewares []middleware
	routes      []route
	exceptions  []exception
}

func (l *list) middleware(task handler) {
	l.middlewares = append(l.middlewares, middleware{task: task})
}

func (l *list) route(method string, pattern string, task handler) {
	l.routes = append(l.routes, route{
		method:  method,
		pattern: pattern,
		task:    task,
	})
}

func (l *list) exception(code string, task handler) {
	l.exceptions = append(l.exceptions, exception{task: task})
}
