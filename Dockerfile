FROM golang:alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

CMD [ "go", "test", "-v", "./..." ]