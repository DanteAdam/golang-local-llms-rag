FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY golang-llms /app/golang-llms
COPY storage /root/storage


RUN chmod +x /app/golang-llms

EXPOSE 8080

CMD ["./golang-llms"]