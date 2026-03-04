ARG GOLANG_VERSION
ARG GOTENBERG_VERSION

FROM golang:$GOLANG_VERSION AS builder

ENV CGO_ENABLED=0

# Define the working directory outside of $GOPATH (we're using go modules).
WORKDIR /home

# Install module dependencies.
COPY go.mod go.sum ./

RUN go mod download &&\
    go mod verify

# Copy the source code.
COPY cmd ./cmd
COPY pkg ./pkg

# Build the binary.
ARG GOTENBERG_VERSION
ARG APP_NAME
ARG APP_VERSION

RUN go build -o gotenberg -ldflags "-X 'github.com/gotenberg/gotenberg/v8/cmd.Version=$GOTENBERG_VERSION-$APP_NAME-$APP_VERSION'" cmd/app/main.go

FROM gotenberg/gotenberg:$GOTENBERG_VERSION

ARG GOTENBERG_VERSION
ARG APP_AUTHOR
ARG APP_REPOSITORY
ARG APP_NAME
ARG APP_VERSION

LABEL org.opencontainers.image.title="Gotenberg $APP_NAME" \
      org.opencontainers.image.description="A containerized API for seamless PDF conversion." \
      org.opencontainers.image.version="$GOTENBERG_VERSION-$APP_NAME-$APP_VERSION" \
      org.opencontainers.image.authors="$APP_AUTHOR" \
      org.opencontainers.image.source="$APP_REPOSITORY"

USER root

# TODO: add your instructions.
ARG GHOSTSCRIPT_VERSION
RUN apt update && apt install -y build-essential git libjpeg-dev libpng-dev libtiff-dev zlib1g-dev libfreetype6-dev libfontconfig1-dev libcupsimage2-dev libgtk-3-dev libidn11-dev libpaper-dev
RUN \
  wget -O ghostscript.tar.gz "https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10060/ghostscript-$GHOSTSCRIPT_VERSION.tar.gz"  &&\
  tar -zxvf ghostscript.tar.gz &&\
  cd ghostscript-$GHOSTSCRIPT_VERSION &&\
  ./configure --without-x &&\
  make -j$(nproc) &&\
  make install

# Copy the Gotenberg binary from the builder stage.
COPY --from=builder /home/gotenberg /usr/bin/

ENV GHOSTSCRIPT_BIN_PATH /usr/local/bin/gs

USER gotenberg