FROM alpine
WORKDIR /app
ENTRYPOINT ["/app/orderly-badger"]
EXPOSE 8080

COPY ./orderly-badger /app/orderly-badger
COPY ./static /app/static