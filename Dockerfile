FROM nikwaf28/alpine-oracle-client:latest

RUN mkdir /app
WORKDIR /app

COPY ./web-generator .
COPY view ./view
COPY assets ./assets

EXPOSE 8082

CMD ["./web-generator"]