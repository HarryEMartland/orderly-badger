FROM alpine
COPY ./orderly-badger /app/orderly-badger
COPY ./static /app/static
WORKDIR /app
ENTRYPOINT ["/app/orderly-badger"]
EXPOSE 8080