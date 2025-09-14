FROM golang:tip-alpine AS builder
COPY . /build
WORKDIR /build
RUN go build -o /snipr . 

FROM alpine:latest AS runner
COPY --from=builder /snipr /snipr
CMD ["/snipr", "--disable_db", "--verbose"]
