FROM alpine:edge

RUN apk --update --no-cache add ca-certificates

COPY /ruble /ruble

CMD ["/ruble"]
