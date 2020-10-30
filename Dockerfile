# build stage
FROM golang:1-alpine as build-env

ADD . /build
WORKDIR /build

RUN go build -o /server

# deploy stage
FROM alpine:latest

EXPOSE 8080

WORKDIR /
COPY --from=build-env /server /
CMD ["/server"]
