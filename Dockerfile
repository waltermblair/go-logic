FROM golang:1.10

# TODO - will passing this in docker-compose override this?
ENV THIS_QUEUE 1

RUN mkdir /app
RUN mkdir -p /go/src/github.com/waltermblair/logic
COPY . /go/src/github.com/waltermblair/logic/
WORKDIR /go/src/github.com/waltermblair/logic

RUN go get -u github.com/golang/dep/...
RUN dep ensure -vendor-only
RUN go build -o /app/main .

ENTRYPOINT ["/app/main"]