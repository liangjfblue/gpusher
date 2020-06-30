FROM golang:1.13 AS build

RUN mkdir /gpusher

WORKDIR	/gpusher

COPY common ./common
COPY proto ./proto
COPY gateway ./gateway
COPY go.mod .

ENV GO11MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR	/gpusher/gateway/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway .



FROM scratch AS prod
COPY --from=build /gpusher/gateway/cmd .
CMD ["./gateway"]
