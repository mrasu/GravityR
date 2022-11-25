FROM golang:1.19
WORKDIR /app
CMD ls

COPY ./mock/go.mod .
COPY ./mock/go.sum .

RUN go mod download

CMD ["go", "run", ".", "serve"]
