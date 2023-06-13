FROM archlinux:latest

LABEL maintainer="yukun@doche.io"
LABEL version="1.0.0"

ENV NAMESPACE="default"
ENV UNIQUEKEY="ctrapi"
ENV IMAGE="docheio/minecraft-be"

RUN pacman-key --init
RUN pacman-key --populate archlinux
RUN pacman -Sy
RUN pacman -S --noconfirm archlinux-keyring
RUN pacman -Syyu --noconfirm
RUN pacman -S --noconfirm go
RUN pacman -S --noconfirm curl
RUN pacman -S --noconfirm tar

WORKDIR /root
RUN curl -sLO https://github.com/docheio/container-api/releases/download/v1.0.0/container-api-v1.0.1.tar.gz
RUN tar zxfp ./container-api-v1.0.0.tar.gz

WORKDIR /root/container-api
RUN go build -o /root/ctrapi

WORKDIR /root
RUN rm -rf container-api-v1.0.0.tar.gz ./container-api
CMD ./ctrapi --namespace ${NAMESPACE} --key ${UNIQUEKEY} --image "${IMAGE}"