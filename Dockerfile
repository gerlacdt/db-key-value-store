FROM alpine:3.6

COPY app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
