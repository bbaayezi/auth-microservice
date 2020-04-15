FROM golang:alpine as builder

WORKDIR /go/src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/auth

FROM scratch

COPY --from=builder /go/bin/auth /go/bin/auth

ENTRYPOINT [ "/go/bin/auth" ]