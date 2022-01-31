FROM golang:latest
WORKDIR /root/
RUN ln -s /usr/bin/dpkg-split /usr/sbin/dpkg-split
RUN ln -s /usr/bin/dpkg-deb /usr/sbin/dpkg-deb
RUN ln -s /bin/rm /usr/sbin/rm
RUN ln -s /bin/tar /usr/sbin/tar
RUN apt-get update && \
      apt-get install -y --no-install-recommends \
      xvfb libfontconfig \
      && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb
RUN apt -f install
RUN go install github.com/ca-gip/hackathon-reward@latest

COPY . .
RUN go mod vendor
RUN GO111MODULE="on" CGO_ENABLED=0 go build -ldflags="-s" -v -o ./build/hackathon-api ./main.go

EXPOSE 8080
CMD ["./build/hackathon-api"]
