FROM golang:1.12.9

ARG cert_path=cert.pem
ARG key_path=secret.key

MAINTAINER Marian Craciunescu

LABEL Description="Meraki health Beats plugin"

RUN \
    apt-get update \
      && apt-get install -y --no-install-recommends \
         netcat \
         python-pip \
         virtualenv \
      && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip
RUN pip install --upgrade setuptools
RUN pip install --upgrade docker-compose==1.23.2

RUN mkdir /plugin
COPY merakibeat.yml /plugin/
COPY fields.yml /plugin/
COPY merakibeat /plugin/merakibeat
#--build-arg key_path=secret.key  cert_path=cert.pem

COPY  $key_path /plugin/server.key
COPY  $cert_path  /plugin/cert.pem
WORKDIR /plugin

ENTRYPOINT ["/plugin/merakibeat", "-e", "-d", "*"]
