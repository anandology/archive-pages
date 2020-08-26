FROM golang:alpine
WORKDIR /app
ADD go.mod /app/go.mod
RUN go mod download
ADD . /app
RUN go build
CMD ["./archive-pages"]
