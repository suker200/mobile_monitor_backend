FROM alpine

RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /

COPY monitor_app /

RUN chmod +x /monitor_app

ENTRYPOINT ["./monitor_app"]