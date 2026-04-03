package admin

import (
	"github.com/go-chi/chi/v5"

	"ekhoes-server/auth"
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

	r.Route("/ctl", func(r chi.Router) {
		r.Get("/sessions", auth.GetSessionsHandler)
		r.Delete("/session/{id}", auth.DeleteSessionHandler)
		r.Delete("/sessions", auth.DeleteAllSessionsHandler)

		r.Get("/connections", websocket.GetConnectionsHandler)

		r.Get("/system", GetSystemInfo)
		r.Get("/top", TopCpuProcesses)
	})

	return nil
}
