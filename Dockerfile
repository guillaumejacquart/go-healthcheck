FROM golang:1.9-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache --update bash git gcc g++ && \
    go get -u -v github.com/kardianos/govendor

WORKDIR /go/src/app
COPY . .

RUN govendor sync && govendor install +local && go build

CMD ["./app"]