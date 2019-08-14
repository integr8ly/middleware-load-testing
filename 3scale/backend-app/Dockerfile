FROM golang:1.12

WORKDIR /usr/src/backend-app

COPY main.go .

RUN go build main.go

EXPOSE 8080

ENTRYPOINT ["./main"]