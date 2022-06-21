# This stage builds the go executable.
FROM golang:1.18.3-buster as go

WORKDIR /root
COPY ./ ./

RUN go build -o bin/app cmd/app/main.go


# Final stage that will be pushed.
FROM debian:buster-slim

FROM node:18.4.0-buster-slim as node

WORKDIR /root

# copy the mermaidcli node package into the container and install
COPY ./mermaidcli/* ./

RUN npm install

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update 2>/dev/null && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		gconf-service \
        libasound2 \
        libatk1.0-0 \
        libatk-bridge2.0-0 \
        libc6 \
        libcairo2 \
        libcups2 \
        libdbus-1-3 \
        libexpat1 \
        libfontconfig1 \
        libgcc1 \
        libgconf-2-4 \
        libgdk-pixbuf2.0-0 \
        libglib2.0-0 \
        libgtk-3-0 \
        libnspr4 \
        libpango-1.0-0 \
        libpangocairo-1.0-0 \
        libstdc++6 \
        libx11-6 \
        libx11-xcb1 \
        libxcb1 \
        libxcomposite1 \
        libxcursor1 \
        libxdamage1 \
        libxext6 \
        libxfixes3 \
        libxi6 \
        libxrandr2 \
        libxrender1 \
        libxss1 \
        libxtst6 \
        libxcb-dri3-0 \
        libgbm1 \
        ca-certificates \
        fonts-liberation \
        libappindicator1 \
        libnss3 \
        lsb-release \
        xdg-utils \
        wget \
        libxshmfence1 \
		2>/dev/null

COPY --from=go /root/bin/app ./app

RUN mkdir -p ./in
RUN mkdir -p ./out
RUN chmod 0777 ./in
RUN chmod 0777 ./out

CMD ["./app", "--mermaid=./node_modules/.bin/mmdc", "--in=./in", "--out=./out", "--puppeteer=./puppeteer-config.json"]

