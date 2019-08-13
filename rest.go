package rest

type App struct {
	list    *list
	api     *API
	handler *Handler
}

func New(prefix string) (*API, *Handler) {
	var l = &list{}
	var api = &API{
		prefix: prefix,
		list:   l,
	}
	var handler = &Handler{
		list: l,
	}
	return api, handler
}
