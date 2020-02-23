FROM golang:1.13.8

ADD . /go/src/application
ADD ./vendor /go/src/application/vendor
WORKDIR /go/src/application

RUN go mod vendor
RUN go build -o /application /go/src/application

CMD ["/application"]