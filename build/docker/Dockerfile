FROM golang:1.17-buster

MAINTAINER dingsqlvine@pingcap.com

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
    sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' \
    /etc/apt/sources.list

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y tzdata locales locales-all
RUN sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && \
    locale-gen
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y file tmux less vim unzip net-tools
RUN apt-get install -y libssl-dev libcurl4-openssl-dev
RUN apt-get install -y git gdb sudo direnv

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GO111MODULE=on
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN echo 'export GOPATH=/go' >> /root/.profile
RUN echo 'export PATH=/usr/local/go/bin:/go/bin:$PATH' >> /root/.profile

RUN go install golang.org/dl/gotip@latest

RUN useradd -ms /bin/bash sqlvine && adduser sqlvine sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN cp /root/.profile /home/sqlvine/.profile
RUN echo 'eval "$(direnv hook bash)"' >> /home/sqlvine/.bashrc

RUN chown sqlvine:sqlvine /go/ -R

USER sqlvine
WORKDIR /home/sqlvine

RUN git clone https://github.com/pingcap/tidb.git
RUN git clone -b dev.fuzz.custom_mutator https://github.com/s3nt3/go.git
RUN cd /home/sqlvine/go/src && ./make.bash
