FROM golang:1.15-alpine AS go
COPY . /web-annie
WORKDIR /web-annie
RUN apk --no-cache add openssl wget tar && wget 'https://github.com/iawia002/lux/releases/download/v0.15.0/lux_0.15.0_Linux_64-bit.tar.gz' && tar -xvf lux_*
RUN apk add --no-cache git 
RUN export GIT_TAG=$(git describe --tags --abbrev=0) && export GIT_HASH=$(git rev-parse HEAD) && export BUILD_DATE=$(date '+%F-%H:%M:%S') && go build -ldflags "-X main.gitTag=$GIT_TAG -X main.gitHash=$GIT_HASH -X main.buildDate=$BUILD_DATE"
				
FROM jrottenberg/ffmpeg:4.0-alpine
COPY --from=go /web-annie /web-annie
WORKDIR /web-annie
ENTRYPOINT ["/web-annie/web-annie"]
