package handler

import (
	"net/http"

	"wildberries_L0/internal/service"
)

type Handler struct {
	service service.Service
	router  *http.ServeMux
}

func (h *Handler) orderPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	order := h.service.GetCache(id)
	if order == nil {
		w.Write([]byte("Такого заказа не существует"))
	}

	w.Write(order)
}

func (h *Handler) Init() {
	h.router.HandleFunc("/order", h.orderPage)

}

func NewHendler(service service.Service, router *http.ServeMux) *Handler {
	return &Handler{
		service: service,
		router:  router,
	}
}
