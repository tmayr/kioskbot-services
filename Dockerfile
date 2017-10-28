FROM golang:alpine

RUN apk update && apk add git
RUN go get github.com/kardianos/govendor
COPY ./ /go/src/kioskbot-services
WORKDIR /go/src/kioskbot-services
RUN govendor fetch +external
RUN go build

EXPOSE 3000
CMD ["./kioskbot-services"]