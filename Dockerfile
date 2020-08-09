FROM golang:1.14-alpine

COPY . /web-annie
WORKDIR /web-annie
RUN ["go", "build"]
ENTRYPOINT ["/web-annie/web-annie"]
