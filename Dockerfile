# ************************************************************************************************************ #
#                                                                                                              #
#                                                      :::::::::  ::::::::   ::::::::   :::    ::: ::::::::::  #
#  Dockerfile                                         :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:          #
#                                                    +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+           #
#  By: ES-Yukun <yukun@doche.io>                    +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#       #
#                                                  +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+             #
#  Created: 2023/08/09 23:48:44 by ES-Yukun       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#              #
#  Updated: 2023/08/09 23:48:46 by ES-Yukun      #########  ########   ########  ###    ### ##########.io.     #
#                                                                                                              #
# ************************************************************************************************************ #

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
RUN curl -sLO https://github.com/docheio/container-api/archive/refs/tags/v1.0.3.tar.gz
RUN tar zxfp ./v1.0.3.tar.gz

RUN mv /root/container-api* /root/container-api
WORKDIR /root/container-api
RUN go build -o /root/ctrapi

WORKDIR /root
RUN rm -rf ./container-api*
CMD ./ctrapi --namespace ${NAMESPACE} --key ${UNIQUEKEY} --image "${IMAGE}"
