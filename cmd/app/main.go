package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tomwright/grace"
	"github.com/tomwright/mermaid-server/internal"
	"os"
)

func main() {
	mermaid := flag.String("mermaid", "", "The full path to the mermaidcli executable.")
	in := flag.String("in", "", "Directory to store input files.")
	out := flag.String("out", "", "Directory to store output files.")
	puppeteer := flag.String("puppeteer", "", "Full path to optional puppeteer config.")
	flag.Parse()

	if *mermaid == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing required argument `mermaid`")
		os.Exit(1)
	}

	if *in == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing required argument `in`")
		os.Exit(1)
	}

	if *out == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing required argument `out`")
		os.Exit(1)
	}

	g := grace.Init(context.Background())

	cache := internal.NewDiagramCache()
	generator := internal.NewGenerator(cache, *mermaid, *in, *out, *puppeteer)

	httpRunner := internal.NewHTTPRunner(generator)
	cleanupRunner := internal.NewCleanupRunner(generator)

	g.Run(httpRunner)
	g.Run(cleanupRunner)

	g.Wait()
}
