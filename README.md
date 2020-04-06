# mermaid-server

Use mermaid-js to generate diagrams in a HTTP endpoint.

While this currently serves the diagrams via HTTP, it could easily be manipulated to server diagrams via other means.

## Basic usage

Start the HTTP server:
```
go run cmd/app/main.go --mermaid=./mermaidcli/node_modules/.bin/mmdc --in=./in --out=./out
```

Send CURL request to generate a diagram:
```
curl --location --request POST 'http://localhost:80/generate' \
--header 'Content-Type: text/plain' \
--data-raw 'graph LR

    A-->B
    B-->C
    C-->D
    C-->F
'
```
