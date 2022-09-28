FROM golang:1.19-alpine3.16

WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build -o http-server

FROM alpine:3.16

WORKDIR /app
COPY --from=0 /app/http-server .
ENTRYPOINT ["/app/http-server"]