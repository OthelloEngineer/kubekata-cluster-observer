FROM golang:1.23

COPY . .
RUN go mod download
EXPOSE 8080

CMD ["go", "run", "main.go"]
