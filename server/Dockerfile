FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY server server/

RUN go build -o neopixels ./server

ENTRYPOINT ["/app/neopixels"]

CMD ["-host=0.0.0.0", "-port=5000"]
