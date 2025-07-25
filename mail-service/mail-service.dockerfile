# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY mailServiceApp /app
COPY templates /templates

CMD [ "/app/mailServiceApp" ]