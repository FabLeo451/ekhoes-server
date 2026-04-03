package admin

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"ekhoes-server/module"
	"ekhoes-server/websocket"
)

var thisModule module.Module

func Register() {
	thisModule = module.Module{
		Id:       "admin",
		Name:     "Admin",
		InitFunc: Init,
	}
	module.Register(thisModule)
}

func Init(r *chi.Mux) error {

	root := fmt.Sprintf("/%s", thisModule.Id)

	r.Route(root, func(r chi.Router) {
		r.Post("/login", Login)

		r.Route("/ctl", func(r chi.Router) {
			r.Get("/sessions", GetSessionsHandler)
			r.Delete("/session/{id}", DeleteSessionHandler)
			r.Delete("/sessions", DeleteAllSessionsHandler)

			r.Get("/connections", websocket.GetConnectionsHandler)

			r.Get("/system", GetSystemInfo)
			r.Get("/top", TopCpuProcesses)
		})
	})

	return nil
}
