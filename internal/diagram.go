package internal

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

// NewDiagram returns a new diagram.
func NewDiagram(description []byte, imgType string) *Diagram {
	return &Diagram{
		description: []byte(strings.TrimSpace(string(description))),
		lastTouched: time.Now(),
		mu:          &sync.RWMutex{},
		imgType:     imgType,
	}
}

// Diagram represents a single diagram.
type Diagram struct {
	// id is the ID of the Diagram
	id string
	// description is the description of the diagram.
	description []byte
	// Output is the filepath to the output file.
	Output string
	// mu is a mutex to protect the last touched value.
	mu *sync.RWMutex
	// lastTouched is the time that the diagram was last used.
	lastTouched time.Time
	// the type of image to generate svg or png
	imgType string
}

// Touch updates the last touched time of the diagram.
func (d *Diagram) Touch() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lastTouched = time.Now()
}

// TouchedInDuration returns true if the diagram has been touched in the given duration.
func (d *Diagram) TouchedInDuration(duration time.Duration) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return time.Now().Add(-duration).Before(d.lastTouched)
}

// ID returns an ID for the diagram.
// The ID is set from the diagram description.
func (d *Diagram) ID() (string, error) {
	if d.id != "" {
		return d.id, nil
	}

	encoded := base64.StdEncoding.EncodeToString(d.description)
	hash := md5.Sum([]byte(encoded))
	d.id = hex.EncodeToString(hash[:]) + d.imgType

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
