# This stage builds the mermaidcli executable.
FROM node:12.12.0-buster as node

WORKDIR /root

# copy the mermaidcli node package into the container and install
COPY ./mermaidcli/* .

RUN npm install


# This stage builds the go executable.
FROM golang:1.13-buster as go

WORKDIR /root
COPY . .

RUN go build -o bin/app cmd/app/main.go


# Final stage that will be pushed.
FROM debian:buster-slim

ENV DEBIAN_FRONTEND=noninteractive
RUN apt update 2>/dev/null && \
	apt install -y --no-install-recommends \
		ca-certificates \
		2>/dev/null

COPY --from=node /root/node_modules/.bin/mmdc ./mermaidcli
COPY --from=go /root/bin/app ./app

# We should now have all the required dependencies to build the proto files.
CMD ["./app", "--mermaid=./mermaidcli"]
