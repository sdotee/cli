FROM golang:1.21 AS builder
LABEL maintainer="mingcheng<mingcheng@outlook.com>"

ENV PACKAGE s.ee/cli
ENV BUILD_DIR ${GOPATH}/src/${PACKAGE}
# ENV GOPROXY https://goproxy.cn,direct

# Build
COPY . ${BUILD_DIR}
WORKDIR ${BUILD_DIR}
RUN make build && cp ./see /bin/see

# Stage2
FROM debian:stable

ENV TZ "Asia/Shanghai"
RUN echo "Asia/Shanghai" > /etc/timezone \
 	&& apt -y update \
	&& apt -y install ca-certificates openssl tzdata curl netcat-openbsd dumb-init \
 	&& apt -y autoremove

COPY --from=builder /bin/see /bin/see

USER nobody
CMD ["/bin/see"]
