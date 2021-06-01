FROM golang:1.16 as builder
WORKDIR /go-secret
ADD . .
RUN make build

FROM alpine:latest
COPY --from=builder /go-secret/secret .
ENTRYPOINT ["./secret"]
CMD ["server"]
