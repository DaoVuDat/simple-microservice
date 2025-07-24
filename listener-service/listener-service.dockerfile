# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY listenerServiceApp /app

CMD [ "/app/listenerServiceApp" ]