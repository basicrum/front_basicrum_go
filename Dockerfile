FROM golang:1.19-alpine as builder

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /go/bin/server .

FROM alpine

RUN apk update
RUN apk add wget

WORKDIR /bin

COPY --from=builder /go/bin /bin

HEALTHCHECK --interval=60s --timeout=10s --retries=5 --start-period=20s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8087/health || exit 1

CMD ["./server"]