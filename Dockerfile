FROM golang AS build

ADD . /entry
WORKDIR /entry

RUN go build -o /go/bin/server wow/cmd/server
RUN go build -o /go/bin/client wow/cmd/client

EXPOSE 5000
