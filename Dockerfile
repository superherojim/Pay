FROM golang:1.20.4-alpine3.16 AS builder

COPY . /data/app
WORKDIR /data/app

RUN rm -rf /data/app/bin/
RUN mkdir -p /data/app/bin/
RUN cd /data/app/cmd/server
RUN export GOPROXY=https://goproxy.cn,direct && cd /data/app/cmd/server && go build -ldflags="-s -w" -o cheemshappypay .
RUN mv /data/app/cmd/server/cheemshappypay /data/app/bin/
RUN mv /data/app/config/prod.yml /data/app/bin/


FROM ubuntu:latest
# 设置环境变量，禁用交互式安装
ENV DEBIAN_FRONTEND=noninteractive

# 将软件源修改为国内的 apt 源
RUN sed -i 's/http:\/\/archive.ubuntu.com\/ubuntu\//http:\/\/mirrors.aliyun.com\/ubuntu\//g' /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    poppler-utils

RUN DEBIAN_FRONTEND=noninteractive apt install -y tzdata

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl

# 设置环境变量
ENV TZ=Asia/Shanghai

ENV APP_ENV=prod

WORKDIR /data/app
COPY --from=builder /data/app/bin /data/app

EXPOSE 8100
ENTRYPOINT [ "./cheemshappypay" ,"-conf","./prod.yml"]

#docker build -t  1.1.1.1:5000/demo-api:v1 --build-arg APP_CONF=config/prod.yml --build-arg  APP_RELATIVE_PATH=./cmd/server/...  .
#docker run -it --rm --entrypoint=ash 1.1.1.1:5000/demo-api:v1