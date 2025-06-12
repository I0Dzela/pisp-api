FROM ubuntu:20.04 as build

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

EXPOSE 7016

# copy
COPY . /usr/local/app

# change working directory
WORKDIR /usr/local/app

# git ssh config
RUN mkdir -p /root/.ssh
COPY docker/.gitconfig /root
COPY docker/config /root/.ssh
COPY docker/secrets/repository-credentials /root/.ssh
RUN chmod 700 /root/.ssh/repository-credentials

RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

# build
RUN make all

# Stage 2:
FROM alpine:3.19.1

# create i0pisp executable directory
RUN mkdir -p /app
RUN mkdir -p /usr/local/bin/i0pisp
RUN mkdir -p /usr/local/bin/i0pisp/openapi

# pisp files binary
COPY --from=build /usr/local/app/bin/ /usr/local/bin/i0pisp

# go binary fix
RUN apk add --no-cache libc6-compat
# Install Apline dependencies
RUN apk add gcompat

# docker container start script
COPY --from=build /usr/local/app/docker/pisp_specs.prod.start.sh /app/start.sh

# copy local ssl certificate
COPY --from=build /usr/local/app/pisp.local.crt /usr/local/bin/i0pisp/pisp.local.crt
COPY --from=build /usr/local/app/pisp.local.key /usr/local/bin/i0pisp/pisp.local.key

COPY --from=build /usr/local/app/rootCA.crt /usr/local/share/ca-certificates/rootCA.crt
RUN cat /usr/local/share/ca-certificates/rootCA.crt >> /etc/ssl/certs/ca-certificates.crt

# openapi specs
COPY --from=build /usr/local/app/openapi/common.yaml /usr/local/bin/i0pisp/openapi/common.yaml
COPY --from=build /usr/local/app/openapi/facekit.yaml /usr/local/bin/i0pisp/openapi/facekit.yaml
COPY --from=build /usr/local/app/openapi/file.yaml /usr/local/bin/i0pisp/openapi/file.yaml
COPY --from=build /usr/local/app/openapi/relation.yaml /usr/local/bin/i0pisp/openapi/relation.yaml
COPY --from=build /usr/local/app/openapi/user.yaml /usr/local/bin/i0pisp/openapi/user.yaml

ENTRYPOINT ["sh", "/app/start.sh"]
