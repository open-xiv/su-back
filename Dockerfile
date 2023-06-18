FROM golang AS builder
WORKDIR $GOPATH/src/su
COPY . .
RUN make build -d

FROM scratch
COPY --from=builder /go/src/su/build/su_amd64 .
COPY web ./web
EXPOSE 8123
CMD ["./su_amd64"]
