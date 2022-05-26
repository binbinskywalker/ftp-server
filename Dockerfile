FROM golang:1.16.4-alpine AS development

ENV PROJECT_PATH=/ftp-server
ENV PATH=$PATH:$PROJECT_PATH/build
ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.13/main/ > /etc/apk/repositories
RUN apk add --no-cache ca-certificates tzdata make git bash protobuf

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements
RUN make

FROM alpine:3.13.2 AS production
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.13/main/ > /etc/apk/repositories
RUN apk --no-cache add ca-certificates tzdata
COPY --from=development /ftp-server/build/ftp-server /usr/bin/ftp-server
USER nobody:nogroup
ENTRYPOINT ["/usr/bin/ftp-server"]
