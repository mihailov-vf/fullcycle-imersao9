FROM golang:1.15 as prod

WORKDIR /go/src

ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt update && \
    apt install build-essential protobuf-compiler librdkafka-dev -y
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    go get google.golang.org/protobuf/cmd/protoc-gen-go && \
    go get github.com/spf13/cobra && \
    wget https://github.com/ktr0731/evans/releases/download/v0.10.8/evans_linux_amd64.tar.gz && \
    tar -xzvf evans_linux_amd64.tar.gz && \
    mv evans /go/bin/ &&  rm -f evans_linux_amd64.tar.gz

CMD [ "tail", "-f", "/dev/null" ]

FROM prod as dev
ARG USERNAME=devuser
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    #
    # [Optional] Add sudo support. Omit if you don't need to install software after connecting.
    && apt-get update \
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    && chown -R $USERNAME:$USERNAME /go

# ********************************************************
# * Anything else you want to do like clean up goes here *
# ********************************************************

# [Optional] Set the default user. Omit if you want to keep the default as root.
USER $USERNAME