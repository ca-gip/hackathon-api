FROM golang:latest as build
WORKDIR $GOPATH/src/github.com/ca-gip/hackathon-api/build
COPY . $GOPATH/src/github.com/ca-gip/hackathon-api/build
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  xvfb libfontconfig wkhtmltopdf \
  && rm -rf /var/lib/apt/lists/*


WORKDIR /root/
COPY --from=build /go/src/github.com/ca-gip/hackathon-api/build/hackathon-api .
EXPOSE 8080
CMD ["./hackathon-api"]
