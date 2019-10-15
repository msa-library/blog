FROM nginx:1.13.12
RUN echo 'Init NGINX'

RUN apt-get update && apt-get install -y certbot

WORKDIR /app

EXPOSE 80 443
