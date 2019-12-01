FROM nginx

COPY build/ /usr/share/nginx/html

RUN mkdir /etc/letsencrypt

RUN rm /etc/nginx/conf.d/default.conf

COPY default.conf /etc/nginx/conf.d/default.conf
