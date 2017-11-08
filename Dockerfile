FROM golang:alpine

ENV PORT 3000

RUN apk update && apk --no-cache add git curl
RUN git clone https://github.com/tmayr/kioskbot-services /go/src/kioskbot-services
RUN go get github.com/kardianos/govendor

WORKDIR /go/src/kioskbot-services
RUN govendor fetch +out
RUN go build

EXPOSE 3000
CMD ["./kioskbot-services"]