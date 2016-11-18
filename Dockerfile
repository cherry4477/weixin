FROM golang:1.6.2

EXPOSE 9090

ENV SERVICE_SOURCE_URL github.com/asiainfoLDP/sever

WORKDIR $GOPATH/src/$SERVICE_SOURCE_URL

ADD . .

RUN go build

CMD ["sh", "-c", "./weixin"]