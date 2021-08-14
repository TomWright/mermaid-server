package internal

import (
	"encoding/json"
	"fmt"
	"github.com/tomwright/grace"
	"github.com/tomwright/gracehttpserverrunner"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NewHTTPRunner returns a grace runner that runs a HTTP server.
func NewHTTPRunner(generator Generator) grace.Runner {
	httpHandler := generateHTTPHandler(generator)

	r := http.NewServeMux()
	r.Handle("/generate", http.HandlerFunc(httpHandler))

	return &gracehttpserverrunner.HTTPServerRunner{
		Server: &http.Server{
			Addr:    ":80",
			Handler: r,
		},
		ShutdownTimeout: time.Second * 5,
	}
}

func writeJSON(rw http.ResponseWriter, value interface{}, status int) {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic("could not marshal value: " + err.Error())
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	if _, err := rw.Write(bytes); err != nil {
		panic("could not write bytes to response: " + err.Error())
	}
}

func writeImage(rw http.ResponseWriter, data []byte, status int, imgType string) error {
	switch imgType {
	case "png":
		rw.Header().Set("Content-Type", "image/png")
	case "svg":
		rw.Header().Set("Content-Type", "image/svg+xml")
	default:
		return fmt.Errorf("unhandled image type: %s", imgType)
	}
	rw.WriteHeader(status)
	if _, err := rw.Write(data); err != nil {
		return fmt.Errorf("could not write image bytes: %w", err)
	}
	return nil
}

func writeErr(rw http.ResponseWriter, err error, status int) {
	log.Printf("[%d] %s", status, err)

	writeJSON(rw, map[string]interface{}{
		"error": err,
	}, status)
}

// URLParam is the URL parameter getDiagramFromGET uses to look for data.
const URLParam = "data"

func getDiagramFromGET(r *http.Request, imgType string) (*Diagram, error) {
	if r.Method != http.MethodGet {
		return nil, fmt.Errorf("expected HTTP method GET")
	}

	queryVal := strings.TrimSpace(r.URL.Query().Get(URLParam))
	if queryVal == "" {
		return nil, fmt.Errorf("missing data")
	}
	data, err := url.QueryUnescape(queryVal)
	if err != nil {
		return nil, fmt.Errorf("could not read query param: %s", err)
	}

	// Create a diagram from the description
	d := NewDiagram([]byte(data), imgType)
	return d, nil
}

func getDiagramFromPOST(r *http.Request, imgType string) (*Diagram, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("expected HTTP method POST")
	}
	// Get description from request body
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %s", err)
	}

	// Create a diagram from the description
	d := NewDiagram(bytes, imgType)
	return d, nil
}

const URLParamImageType = "type"

// generateHTTPHandler returns a HTTP handler used to generate a diagram.
func generateHTTPHandler(generator Generator) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var diagram *Diagram

		imgType := r.URL.Query().Get(URLParamImageType)

		switch imgType {
		case "png", "svg":
		case "":
			imgType = "svg"
		default:
			writeErr(rw, fmt.Errorf("unsupported image type (%s) use svg or png", imgType), http.StatusBadRequest)
			return
		}

		var err error
		switch r.Method {
		case http.MethodGet:
			diagram, err = getDiagramFromGET(r, imgType)
		case http.MethodPost:
			diagram, err = getDiagramFromPOST(r, imgType)
		default:
			writeErr(rw, fmt.Errorf("unexpected HTTP method %s", r.Method), http.StatusBadRequest)
			return
		}
		if err != nil {
			writeErr(rw, err, http.StatusBadRequest)
			return
		}
		if diagram == nil {
			writeErr(rw, fmt.Errorf("could not create diagram"), http.StatusInternalServerError)
			return
		}

		// Generate the diagram
		if err := generator.Generate(diagram); err != nil {
			writeErr(rw, fmt.Errorf("could not generate diagram: %s", err), http.StatusInternalServerError)
			return
		}

		// Output the diagram as an SVG.
		// We assume generate always generates an SVG at this point in time.
		diagramBytes, err := ioutil.ReadFile(diagram.Output)
		if err != nil {
			writeErr(rw, fmt.Errorf("could not read diagram bytes: %s", err), http.StatusInternalServerError)
			return
		}
		if err := writeImage(rw, diagramBytes, http.StatusOK, imgType); err != nil {
			writeErr(rw, fmt.Errorf("could not write diagram: %w", err), http.StatusInternalServerError)
		}
	}
}
