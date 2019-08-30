package rest

func New(prefix string) API {
	var _list = new(list)

	var _handler = &handler{
		list: _list,
	}

	var api = &api{
		prefix:  trim(prefix),
		list:    _list,
		handler: _handler,
	}

	return api
}
