### Default image is base. You can add other support by modifying BASE_IMAGE_TAG. The following parameters are supported: base (default), aria2, ffmpeg, aio
ARG BASE_IMAGE_TAG=ffmpeg

FROM chwjbn/apline:prod AS builder
LABEL stage=go-builder

# 替换为阿里云 apk 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache bash curl jq gcc git go musl-dev

WORKDIR /app/
COPY ./src ./

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# ENV HTTP_PROXY=http://192.168.5.75:54109
# ENV HTTPS_PROXY=http://192.168.5.75:54109
RUN bash build.sh dev docker

FROM openlistteam/openlist-base-image:${BASE_IMAGE_TAG}
LABEL MAINTAINER="OpenList"
ARG INSTALL_FFMPEG=true
ARG INSTALL_ARIA2=false
ARG USER=openlist
ARG UID=1001
ARG GID=1001

WORKDIR /opt/openlist/

RUN addgroup -g ${GID} ${USER} && \
    adduser -D -u ${UID} -G ${USER} ${USER} && \
    mkdir -p /opt/openlist/data

COPY --from=builder --chmod=755 --chown=0:0 /app/bin/openlist ./
COPY --from=builder --chmod=755 --chown=0:0 /app/entrypoint.sh /entrypoint.sh

USER 0
RUN /entrypoint.sh version

ENV PUID=0 PGID=0 UMASK=022 RUN_ARIA2=${INSTALL_ARIA2}
VOLUME /opt/openlist/data/
EXPOSE 5244 5245
CMD [ "/entrypoint.sh" ]
