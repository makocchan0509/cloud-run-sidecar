FROM golang:1.20-alpine

WORKDIR /go/src/

COPY . ./
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build main.go

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./main"]