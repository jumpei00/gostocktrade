FROM golang:1.16.3

WORKDIR /go/src/gostocktrade
COPY go.mod go.sum ./
RUN go mod tidy

CMD [ "go", "run", "main.go" ]