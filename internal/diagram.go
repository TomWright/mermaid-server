package internal

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// NewDiagram returns a new diagram.
func NewDiagram(description []byte) *Diagram {
	return &Diagram{
		description: []byte(strings.TrimSpace(string(description))),
	}
}

type Diagram struct {
	// iD is the ID of the Diagram
	id string
	// description is the description of the diagram.
	description []byte
	// Output is the filepath to the output file.
	Output string
}

// ID returns an ID for the diagram.
// The ID is set from the diagram description.
func (d *Diagram) ID() (string, error) {
	if d.id != "" {
		return d.id, nil
	}

	encoded := base64.StdEncoding.EncodeToString(d.description)
	hash := md5.Sum([]byte(encoded))
	d.id = hex.EncodeToString(hash[:])

	return d.id, nil
}

// Description returns the diagram description.
func (d *Diagram) Description() []byte {
	return d.description
}

// Description returns the diagram description.
func (d *Diagram) WithDescription(description []byte) *Diagram {
	d.description = description
	d.id = ""
	return d
}
