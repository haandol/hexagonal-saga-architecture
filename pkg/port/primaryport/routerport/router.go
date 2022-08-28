package routerport

type RouterGroup interface {
	Use(middleware ...any)
	Group(path string) RouterGroup
	Handle(method string, relativePath string, handlerFunc ...any)
}

type Router interface {
	Route(RouterGroup)
}
