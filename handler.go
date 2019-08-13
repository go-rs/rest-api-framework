package rest

import (
	"fmt"
	"net/http"
)

type Handler struct {
	list *list
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, handle := range h.list.routes {
		fmt.Println(handle)
		handle.task()
	}
	_, _ = w.Write([]byte("Hello, World!"))
}
