FROM golang:1.19-alpine

WORKDIR /build

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /api cmd/api/main.go

ENTRYPOINT ["/api"]
