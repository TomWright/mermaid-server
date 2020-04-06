package main

import (
	"flag"
	"fmt"
	"mermaid-server/internal"
	"net/http"
	"os"
)

var testInput = []byte(`
graph TB

    subgraph "Jira"
        createTicket["Create ticket"]
        updateTicket["Update ticket"]
        fireWebhook["Fire webhook"]

        createTicket-->fireWebhook
        updateTicket-->fireWebhook
    end

    subgraph "Jira Webhook"
        receiveWebhook["Receive webhook"]
        storeEvent["Store event in immutable list"]
        publishNewStoredEventEvent["Publish message to notify system of new event"]

        createEvent["Create internal event that will be stored"]
        setCreatedAt["Set created at date to now"]
        setSourceJira["Set source to Jira webhook"]

        receiveWebhook-->createEvent-->setCreatedAt-->setSourceJira-->storeEvent
        storeEvent-->publishNewStoredEventEvent
    end

    fireWebhook-->receiveWebhook

    subgraph "Play Event"
        publishEventUpdated["Publish message to notify system of new status"]

        verifyEventSource["Verify event source"]
        parsePayload["Parse event payload using source to determine structure"]
        findEventHandler["Find the handler for the specific event type + version"]
        getLatestPersistedState["Get latest persisted state"]
        changeInMemoryStateUsingEventData["Change in-memory state using event data"]
        persistUpdatedState["Persist updated state"]

        verifyEventSource-->parsePayload
        parsePayload-->findEventHandler
        findEventHandler-->getLatestPersistedState-->changeInMemoryStateUsingEventData-->persistUpdatedState

        persistUpdatedState-->publishEventUpdated
    end

    publishNewStoredEventEvent-->verifyEventSource
`)

func main() {
	mermaid := flag.String("mermaid", "", "The full path to the mermaidcli executable.")
	in := flag.String("in", "", "Directory to store input files.")
	out := flag.String("out", "", "Directory to store output files.")
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

	cache := internal.NewDiagramCache()
	generator := internal.NewGenerator(cache, *mermaid, *in, *out)

	httpHandler := internal.GenerateHTTPHandler(generator)

	r := http.NewServeMux()
	r.Handle("/generate", http.HandlerFunc(httpHandler))

	httpServer := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	_, _ = fmt.Fprintf(os.Stdout, "Listening on address %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		_, _ = fmt.Fprintf(os.Stderr, "Could not listen for http connections: %s", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Shutdown")
}
