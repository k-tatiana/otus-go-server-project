FROM golang:1.24
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

EXPOSE 8080
CMD ["./main"]