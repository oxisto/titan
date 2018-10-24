FROM node AS build-frontend

WORKDIR /tmp

ADD frontend/*.json ./
ADD frontend/*.lock ./
RUN yarn install

ADD frontend/. .
RUN yarn run lint
RUN yarn run build --prod

FROM golang:stretch AS build-server

WORKDIR /go/src/github.com/oxisto/titan/server

RUN apt update && apt -y install unzip

# copy SDE version and download EVE SDE
COPY sde.* ./
RUN ./sde.sh

# install dep utility
RUN go get -u github.com/golang/dep/cmd/dep

# copy dependency information and fetch them
COPY Gopkg.* ./
RUN dep ensure --vendor-only

# copy sources
COPY . .

# build and install (without C-support, otherwise there issue because of the musl glibc replacement on Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -a cmd/server/server.go

FROM debian:9-slim
# update CA certificates
RUN apt update && apt install -y ca-certificates
WORKDIR /usr/titan
COPY --from=build-frontend /tmp/dist ./frontend/dist
COPY --from=build-server /go/src/github.com/oxisto/titan/server .
COPY --from=build-server /go/src/github.com/oxisto/titan/sde.version .
COPY --from=build-server /go/src/github.com/oxisto/titan/sde/ sde
CMD ["./server"]
