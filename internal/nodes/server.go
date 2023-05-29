package nodes

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler, shutdown chan struct{}) {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go s.httpServer.ListenAndServe()

	select {
	case <-shutdown:
		s.httpServer.Close()
	}
}
