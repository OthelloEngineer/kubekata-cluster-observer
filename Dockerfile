FROM golang:1.23-bookworm

COPY . .
RUN go mod download
EXPOSE 8081
CMD ["go", "run", "main.go"]
