package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

// Generator provides the ability to generate a diagram.
type Generator interface {
	// Generate generates the given diagram.
	Generate(diagram *Diagram) error
	// CleanUp removes any diagrams that haven't used within the given duration.
	CleanUp(duration time.Duration) error
}

// NewGenerator returns a generator that can be used to generate diagrams.
func NewGenerator(cache DiagramCache, mermaidCLIPath string, inPath string, outPath string, puppeteerConfigPath string) Generator {
	return &cachingGenerator{
		cache:               cache,
		mermaidCLIPath:      mermaidCLIPath,
		inPath:              inPath,
		outPath:             outPath,
		puppeteerConfigPath: puppeteerConfigPath,
	}
}

// cachingGenerator is an implementation of Generator.
type cachingGenerator struct {
	cache               DiagramCache
	mermaidCLIPath      string
	inPath              string
	outPath             string
	puppeteerConfigPath string
}

// Generate generates the given diagram.
func (c cachingGenerator) Generate(diagram *Diagram) error {
	has, err := c.cache.Has(diagram)
	if err != nil {
		return fmt.Errorf("cache.Has failed: %w", err)
	}
	if has {
		cached, err := c.cache.Get(diagram)
		if err != nil {
			return fmt.Errorf("cache.Get failed: %w", err)
		}
		*diagram = *cached

		// Update diagram last touched date
		diagram.Touch()
		if err := c.cache.Store(diagram); err != nil {
			return fmt.Errorf("cache.Store failed: %w", err)
		}

		return nil
	}

	diagram.Touch()
	if err := c.generate(diagram); err != nil {
		return fmt.Errorf("cachingGenerater.generate failed: %w", err)
	}
	if err := c.cache.Store(diagram); err != nil {
		return fmt.Errorf("cache.Store failed: %w", err)
	}
	return nil
}

// generate does the actual file generation.
func (c cachingGenerator) generate(diagram *Diagram) error {
	id, err := diagram.ID()
	if err != nil {
		return fmt.Errorf("cannot get diagram ID: %w", err)
	}

	inPath := fmt.Sprintf("%s/%s.mmd", c.inPath, id)
	outPath := fmt.Sprintf("%s/%s.%s", c.outPath, id, diagram.imgType)

	if err := ioutil.WriteFile(inPath, diagram.description, 0644); err != nil {
		return fmt.Errorf("could not write to input file [%s]: %w", inPath, err)
	}

	_, err = os.Stat(c.mermaidCLIPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("mermaid executable does not exist: %w", err)
	}
	if err != nil {
		return fmt.Errorf("could not stat mermaid executable: %w", err)
	}

	args := []string{
		"-i", inPath,
		"-o", outPath,
	}
	if c.puppeteerConfigPath != "" {
		args = append(args, "-p", c.puppeteerConfigPath)
	}

	cmd := exec.Command(c.mermaidCLIPath, args...)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = bufio.NewWriter(&stdOut)
	cmd.Stderr = bufio.NewWriter(&stdErr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed when executing mermaid: %w: %s: %s", err, string(stdOut.Bytes()), string(stdErr.Bytes()))
	}
	log.Printf("Generated: %s: %s: %s", id, string(stdOut.Bytes()), string(stdErr.Bytes()))

	diagram.Output = outPath

	return nil
}

// CleanUp removes any diagrams that haven't used within the given duration.
func (c cachingGenerator) CleanUp(duration time.Duration) error {
	log.Printf("Running cleanup")
	diagrams, err := c.cache.GetAll()
	if err != nil {
		return fmt.Errorf("could not get cached diagrams: %w", err)
	}
	for _, d := range diagrams {
		if !d.TouchedInDuration(duration) {
			if err := c.delete(d); err != nil {
				return fmt.Errorf("could not delete diagram: %w", err)
			}
		}
	}
	return nil
}

// delete removes any diagrams that haven't used within the given duration.
func (c cachingGenerator) delete(diagram *Diagram) error {
	id, err := diagram.ID()
	if err != nil {
		return fmt.Errorf("cannot get diagram ID: %w", err)
	}

	log.Printf("Cleaning up diagram: %s", id)

	inPath := fmt.Sprintf("%s/%s.mmd", c.inPath, id)
	outPath := fmt.Sprintf("%s/%s.svg", c.outPath, id)

	if err := os.Remove(inPath); err != nil {
		return fmt.Errorf("could not delete diagram input: %w", err)
	}
	if err := os.Remove(outPath); err != nil {
		return fmt.Errorf("could not delete diagram output: %w", err)
	}
	if err := c.cache.Delete(diagram); err != nil {
		return fmt.Errorf("could not remove diagram from cache: %w", err)
	}

	return nil
}
