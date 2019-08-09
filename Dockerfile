FROM golang:1.12 as builder

WORKDIR /speedtest_exporter

ADD cmd cmd
ADD pkg pkg
ADD internal internal
ADD go.mod go.sum ./

RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/app cmd/speedtest_exporter/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/app /bin/

ENTRYPOINT ["/bin/app"]
