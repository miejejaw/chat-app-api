FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/app

RUN go build -o /cmd/app/main .

EXPOSE 8080

#COPY .env .env

CMD ["/cmd/app/main"]
