FROM golang:1.16-alpine AS build
LABEL authors="lBAH, Kalinin Iwan <koefic.cien@gmail.com>"

COPY server.go .
COPY client.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server server.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/client client.go

FROM ghcr.io/1bah/fenrir-base:latest AS fenrir

FROM ubuntu:latest
LABEL authors="lBAH, Kalinin Iwan <koefic.cien@gmail.com>"

RUN apt update && apt install --no-install-recommends -y \
    sudo \
    systemd \
    ntfs-3g \
    udev \
    libbpf1 \
    gpg \
    && apt clean

RUN useradd -G sudo -m -s "/bin/bash" fikus && \
    groupadd --system reckod

WORKDIR /home/fikus

COPY --from=fenrir /usr/local/bin/fenrir* /usr/local/bin/
COPY --from=fenrir /usr/local/bin/loki /usr/local/bin/
COPY --from=fenrir --chown=fikus:fikus /etc/fenrir /etc/fenrir
COPY --from=fenrir --chown=fikus:fikus /home/fikus/.fenrir /home/fikus/.fenrir

RUN echo "fikus:breach" | chpasswd && \
    echo "fikus ALL=(ALL) ALL" >> /etc/sudoers

COPY --from=build /bin/server /usr/bin/server
COPY --from=build /bin/client /usr/bin/client

COPY wsl.fs /wsl.fs
COPY entrypoint /entrypoint
COPY pubkey.gpg /key.gpg

ENTRYPOINT [ "/entrypoint" ]
