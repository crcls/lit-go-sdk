FROM curlimages/curl:latest

WORKDIR /tmp
RUN curl -L https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz | tar xvz

FROM golang:1.18.4

COPY --from=0 /tmp/modd-0.8-linux64/modd /usr/local/bin

RUN mkdir /tmp/cache
ENV GOMODCACHE=/tmp/cache

WORKDIR /go/src

CMD ["modd"]
