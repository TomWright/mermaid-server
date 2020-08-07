package internal

import (
	"net/http"
)

// NewHTTPService returns a service that can be used to start a http server
// that will generate diagrams.
func NewHTTPService(generator Generator) *httpService {
	return &httpService{
		generator: generator,
	}
}

// httpService is a service that can be used to start a http server
// that will generate diagrams.
type httpService struct {
	httpServer *http.Server
	generator  Generator
}

// Start starts the HTTP server.
func (s *httpService) Start() error {
	httpHandler := generateHTTPHandler(s.generator)

	r := http.NewServeMux()
	r.Handle("/generate", http.HandlerFunc(httpHandler))

	s.httpServer = &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func (s *httpService) Stop() {
	if s != nil {
		_ = s.httpServer.Close()
	}
}
