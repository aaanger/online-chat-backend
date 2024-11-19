package cmd

import "net/http"

type Server struct {
	httpServer *http.Server
}

func (srv *Server) Run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:    port,
		Handler: handler,
	}
	return srv.httpServer.ListenAndServe()
}
