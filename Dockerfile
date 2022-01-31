FROM golang:latest
ARG TARGETOS TARGETARCH
WORKDIR /root/
RUN ln -s /usr/bin/dpkg-split /usr/sbin/dpkg-split
RUN ln -s /usr/bin/dpkg-deb /usr/sbin/dpkg-deb
RUN ln -s /bin/rm /usr/sbin/rm
RUN ln -s /bin/tar /usr/sbin/tar

COPY . .
RUN go mod vendor
RUN GO111MODULE="on" GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags="-s" -v -o ./build/hackathon-api ./main.go


FROM ubuntu:20.04
ARG TARGETARCH
WORKDIR /app
COPY --from=0 /root/build/hackathon-api /app/hackathon-api

RUN apt update && \
      apt install -y --no-install-recommends \
      xvfb libfontconfig wget fontconfig xfonts-75dpi xfonts-100dpi xfonts-scalable xfonts-base \
      && rm -rf /var/lib/apt/lists/*

RUN wget --no-check-certificate https://github.com/ca-gip/hackathon-api/releases/download/v0.1.0/hackathon-reward-${TARGETARCH} && \
    chmod a+x hackathon-reward-${TARGETARCH} && \
    mv ./hackathon-reward-${TARGETARCH} /usr/local/bin/hackathon-reward

RUN wget --no-check-certificate https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_${TARGETARCH}.deb
RUN apt update && \
    apt install ./wkhtmltox_0.12.6-1.focal_${TARGETARCH}.deb  -y && \
    rm wkhtmltox_0.12.6-1.focal_${TARGETARCH}.deb


EXPOSE 8080
CMD ["/app/hackathon-api"]
