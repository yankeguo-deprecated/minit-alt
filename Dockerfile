FROM golang:1.14 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go test -mod vendor -v
RUN go build -mod vendor -o /minit

FROM alpine:3.12
COPY --from=builder /minit /minit
CMD ["/minit"]