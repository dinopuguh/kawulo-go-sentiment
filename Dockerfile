FROM golang:1.14.1

RUN mkdir -p /go/src/github.com/dinopuguh/kawulo-sentiment/

WORKDIR /go/src/github.com/dinopuguh/kawulo-sentiment/

COPY . .

RUN go build -o sentiment main.go

EXPOSE 9090

CMD /go/src/github.com/dinopuguh/kawulo-sentiment/sentiment