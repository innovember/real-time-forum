FROM golang:latest
LABEL maintainer="github.com/innovember"
ENV GO111MODULE=on
WORKDIR /go/src/real-time-forum

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

EXPOSE 8081

RUN go build -o ./build/api ./cmd/api/main.go

CMD ["./build/api"]