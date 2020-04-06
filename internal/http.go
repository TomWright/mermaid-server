package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	writeJSON(rw, map[string]interface{}{
		"error": err,
	}, status)
}

func GenerateHTTPHandler(generator Generator) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(rw, fmt.Errorf("expected HTTP method POST"), http.StatusBadRequest)
			return
		}
		// Get description from request body
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeErr(rw, fmt.Errorf("could not read body: %s", err), http.StatusInternalServerError)
			return
		}

		// Create a diagram from the description
		d := NewDiagram(bytes)

		// Generate the diagram
		if err := generator.Generate(d); err != nil {
			writeErr(rw, fmt.Errorf("could not generate diagram: %s", err), http.StatusInternalServerError)
			return
		}

		// Output the diagram as an SVG.
		// We assume generate always generates an SVG at this point in time.
		diagramBytes, err := ioutil.ReadFile(d.Output)
		if err != nil {
			writeErr(rw, fmt.Errorf("could not read diagram bytes: %s", err), http.StatusInternalServerError)
			return
		}
		writeSVG(rw, diagramBytes, http.StatusOK)
	}
}
