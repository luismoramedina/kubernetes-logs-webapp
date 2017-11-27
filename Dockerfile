FROM google/debian:wheezy
MAINTAINER Luis Mora Medina <luismoramedina@gmail.com>

ADD klogs klogs
ADD main.html main.html
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/klogs"]