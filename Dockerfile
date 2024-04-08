FROM golang:1.21-alpine3.19 as builder
RUN mkdir /short_url
RUN mkdir /short_url/short
RUN mkdir /short_url/short/cache
RUN mkdir /short_url/short/model
RUN mkdir /short_url/short/repository
RUN mkdir /short_url/builder 
RUN mkdir /short_url/handler 

COPY go.mod /short_url/
COPY go.sum /short_url/
COPY main.go /short_url/
COPY /short/cache/cache.go /short_url/short/cache/
COPY /short/model/model.go /short_url/short/model/
COPY /short/repository/repository.go /short_url/short/repository/
COPY /builder/builder.go /short_url/builder/
COPY /handler/handler.go /short_url/handler/

WORKDIR /short_url
RUN cd /short_url
RUN go build -ldflags="-extldflags=-static" -o shorturl /short_url/main.go

FROM alpine:latest
COPY --from=builder /short_url/shorturl .
ENTRYPOINT ["/shorturl"]
