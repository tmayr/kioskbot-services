FROM golang:alpine

RUN go install github.com/tmayr/kioskbot-services

EXPOSE 3000
CMD ["/go/bin/kioskbot-services"]