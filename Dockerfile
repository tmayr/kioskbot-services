FROM golang:alpine

RUN apk update && apk --no-cache add git curl
RUN git clone https://github.com/tmayr/kioskbot-services /go/src/kioskbot-services
RUN go get github.com/kardianos/govendor

WORKDIR /go/src/kioskbot-services
RUN govendor fetch +out
RUN go build

CMD ["./kioskbot-services"]