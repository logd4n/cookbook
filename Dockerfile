FROM golang:1.22.4-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o cookbook ./cmd/main.go
RUN chmod +x cookbook
CMD ["./cookbook"]