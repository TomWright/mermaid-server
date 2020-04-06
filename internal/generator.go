package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
)

// Generator provides the ability to generate a diagram.
type Generator interface {
	// Generate generates the given diagram.
	Generate(diagram *Diagram) error
}

func NewGenerator(cache DiagramCache, mermaidCLIPath string, inPath string, outPath string) Generator {
	return &cachingGenerator{
		cache:          cache,
		mermaidCLIPath: mermaidCLIPath,
		inPath:         inPath,
		outPath:        outPath,
	}
}

// cachingGenerator is an implementation of Generator.
type cachingGenerator struct {
	cache          DiagramCache
	mermaidCLIPath string
	inPath         string
	outPath        string
}

// Generate generates the given diagram.
func (c cachingGenerator) Generate(diagram *Diagram) error {
	has, err := c.cache.Has(diagram)
	if err != nil {
		return err
	}
	if has {
		cached, err := c.cache.Get(diagram)
		if err != nil {
			return err
		}
		*diagram = *cached
		return nil
	}
	if err := c.generate(diagram); err != nil {
		return err
	}
	if err := c.cache.Store(diagram); err != nil {
		return err
	}
	return nil
}

// generate does the actual file generation.
func (c cachingGenerator) generate(diagram *Diagram) error {
	id, err := diagram.ID()
	if err != nil {
		return err
	}

	has, err := c.cache.Has(diagram)
	if err != nil {
		return err
	}
	if has {
		cached, err := c.cache.Get(diagram)
		if err != nil {
			return err
		}
		*diagram = *cached
		return nil
	}

	inPath := fmt.Sprintf("%s/%s.mmd", c.inPath, id)
	outPath := fmt.Sprintf("%s/%s.svg", c.outPath, id)

	if err := ioutil.WriteFile(inPath, diagram.description, 0644); err != nil {
		return err
	}

	cmd := exec.Command(c.mermaidCLIPath, "-i", inPath, "-o", outPath)
	var stdOut bytes.Buffer
	cmd.Stdout = bufio.NewWriter(&stdOut)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, string(stdOut.Bytes()))
	}

	diagram.Output = outPath

	return nil
}
