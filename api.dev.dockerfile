FROM ubuntu:20.04

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

ENV TZ=Europe/Sarajevo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt update
RUN apt upgrade -y
RUN apt dist-upgrade -y
RUN apt install -y tzdata
RUN apt install curl htop git zip nano build-essential pkg-config \
  vim lsof wget -y

RUN apt update

# Go
RUN wget -c https://dl.google.com/go/go1.23.4.linux-amd64.tar.gz -O - | tar -xz -C /usr/local
ENV PATH="/usr/local/go/bin:$PATH"

RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
RUN go install github.com/golang/mock/mockgen@v1.6.0

ENV GO_PATH=/root/go
ENV PATH="$PATH:$GO_PATH/bin"

RUN apt update

EXPOSE 6030

# change working directory
WORKDIR /app

# git ssh config
RUN mkdir -p /root/.ssh
COPY docker/.gitconfig /root
COPY docker/config /root/.ssh
COPY docker/secrets/repository-credentials /root/.ssh
RUN chmod 700 /root/.ssh/repository-credentials

RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

RUN mkdir -p /usr/local/share/ca-certificates/pisp
COPY rootCA.crt /usr/local/share/ca-certificates/pisp/rootCA.crt
RUN update-ca-certificates

# SwaggerHub CLI
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
RUN apt install -y nodejs
RUN npm i -g swaggerhub-cli