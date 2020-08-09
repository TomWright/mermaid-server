package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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

func writeSVG(rw http.ResponseWriter, data []byte, status int) {
	rw.Header().Set("Content-Type", "image/svg+xml")
	rw.WriteHeader(status)
	if _, err := rw.Write(data); err != nil {
		panic("could not write bytes to response: " + err.Error())
	}
}

func writeErr(rw http.ResponseWriter, err error, status int) {
	log.Printf("[%d] %s", status, err)

	writeJSON(rw, map[string]interface{}{
		"error": err,
	}, status)
}

// URLParam is the URL parameter getDiagramFromGET uses to look for data.
const URLParam = "data"

func getDiagramFromGET(rw http.ResponseWriter, r *http.Request) *Diagram {
	if r.Method != http.MethodGet {
		writeErr(rw, fmt.Errorf("expected HTTP method GET"), http.StatusBadRequest)
		return nil
	}

	queryVal := strings.TrimSpace(r.URL.Query().Get(URLParam))
	if queryVal == "" {
		writeErr(rw, fmt.Errorf("missing data"), http.StatusBadRequest)
		return nil
	}
	data, err := url.QueryUnescape(queryVal)
	if err != nil {
		writeErr(rw, fmt.Errorf("could not read query param: %s", err), http.StatusBadRequest)
		return nil
	}

	// Create a diagram from the description
	d := NewDiagram([]byte(data))
	return d
}

func getDiagramFromPOST(rw http.ResponseWriter, r *http.Request) *Diagram {
	if r.Method != http.MethodPost {
		writeErr(rw, fmt.Errorf("expected HTTP method POST"), http.StatusBadRequest)
		return nil
	}
	// Get description from request body
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErr(rw, fmt.Errorf("could not read body: %s", err), http.StatusInternalServerError)
		return nil
	}

	// Create a diagram from the description
	d := NewDiagram(bytes)
	return d
}

// generateHTTPHandler returns a HTTP handler used to generate a diagram.
func generateHTTPHandler(generator Generator) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var diagram *Diagram

		switch r.Method {
		case http.MethodGet:
			diagram = getDiagramFromGET(rw, r)
		case http.MethodPost:
			diagram = getDiagramFromPOST(rw, r)
		default:
			writeErr(rw, fmt.Errorf("unexpected HTTP method %s", r.Method), http.StatusBadRequest)
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
		writeSVG(rw, diagramBytes, http.StatusOK)
	}
}
