FROM golang:1.20.3-alpine as build-env

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o /app/build .

EXPOSE 8000
ENTRYPOINT [ "/app/build" ]