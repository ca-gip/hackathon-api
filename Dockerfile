FROM golang:latest
ARG TARGETARCH TARGETOS

RUN echo "I am running on $TARGETARCH building for $TARGETPLATFORM" > /log

WORKDIR /root/
RUN ln -s /usr/bin/dpkg-split /usr/sbin/dpkg-split
RUN ln -s /usr/bin/dpkg-deb /usr/sbin/dpkg-deb
RUN ln -s /bin/rm /usr/sbin/rm
RUN ln -s /bin/tar /usr/sbin/tar

COPY . .
RUN go mod vendor
RUN GO111MODULE="on" CGO_ENABLED=0 go build  -a -ldflags '-extldflags "-static"' -v -o ./build/hackathon-api ./main.go

RUN apt update && \
      apt install -y --no-install-recommends \
      xvfb libfontconfig wget fontconfig xfonts-75dpi xfonts-100dpi xfonts-scalable xfonts-base \
      && rm -rf /var/lib/apt/lists/*


RUN wget  https://github.com/ca-gip/hackathon-api/releases/download/v0.1.1/hackathon-reward-${TARGETOS}-${TARGETARCH} -O hackathon-reward
RUN chmod a+x hackathon-reward && \
    mv ./hackathon-reward /usr/local/bin/hackathon-reward

RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_${TARGETARCH}.deb
RUN apt update && \
    apt install ./wkhtmltox_0.12.6-1.buster_${TARGETARCH}.deb  -y && \
    rm wkhtmltox_0.12.6-1.buster_${TARGETARCH}.deb


EXPOSE 8080
CMD ["./build/hackathon-api"]
