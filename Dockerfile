FROM alpine:3.19.1

MAINTAINER zpz zpz@qq.com

RUN apk --update add redis && apk add --no-cache libc6-compat

ADD ./redisKeySample /data/soft/redisKeySample

# 设置时区
RUN rm -f /etc/localtime \
&& ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

# 设置启动访问的初始位置，即工作目录，登录落脚点
ENV MYPATH /data/soft/
WORKDIR $MYPATH