FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM golang:1.16 AS builder
WORKDIR /gomod/noolite
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o /go/bin/noolite ./cmd

FROM scratch
# configurations
WORKDIR /app
# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# the main program:
COPY --from=builder /go/bin/noolite ./noolite
CMD ["./noolite"]
