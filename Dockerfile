FROM golang:1.20-alpine as builder

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /go/bin/server .

FROM alpine

RUN apk update
RUN apk add wget bash

WORKDIR /bin

COPY --from=builder /go/bin /bin

CMD ["./server"]