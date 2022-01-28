FROM golang:latest
WORKDIR /root/
COPY build/hackathon-api .
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  xvfb libfontconfig wkhtmltopdf \
  && rm -rf /var/lib/apt/lists/*

EXPOSE 8080
CMD ["./hackathon-api"]
