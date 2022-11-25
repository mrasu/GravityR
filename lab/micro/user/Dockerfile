FROM golang:1.19

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

CMD ["go", "run", "."]
