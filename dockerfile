FROM golang AS compiling_stage
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
ADD main.go .
ADD go.mod .
ADD ./src/ /go/src/pipeline/src
RUN go install .

FROM alpine:latest
LABEL version="1.0.1"
LABEL maintainer="Artem Rybakov<rybakov333@gmail.com>"
WORKDIR /root/
COPY --from=compiling_stage /go/bin/pipeline .
ENTRYPOINT ./pipeline