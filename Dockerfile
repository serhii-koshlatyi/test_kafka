FROM golang:1.14-alpine AS build

ENV GO111MODULE=on

WORKDIR /go/src/app

RUN apk add bash make

COPY . .
RUN make deps

RUN export BLDDIR=/go/bin && \
    make clean && \
    make build

RUN ls -al /go/bin/

FROM alpine

WORKDIR /app

COPY ./config ./config

COPY --from=build /go/bin/ ./
CMD ["./test"]