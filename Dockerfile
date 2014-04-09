# Docker environment for httpforwarder
#
# There are helper scripts in the `docker` subdirectory of this repo.
#
# To hack on httpforwarder, build & run with the codebase as a linked volume:
#
#  $ docker/build
#  $ docker/run
#
# To publish a release to the Docker registry:
#
#  $ docker/release

FROM       quay.io/99designs/base:20140402-193006
MAINTAINER michael.tibben@99designs.com

# Update apt index; base image may be stale, or from a different mirror.
RUN apt-get update

# Golang
RUN curl -s https://go.googlecode.com/files/go1.2.1.src.tar.gz | tar -v -C /usr/local -xz
ENV PATH   /usr/local/go/bin:$PATH
ENV GOPATH /go:/app
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Upload app source. This may be masked using volumes in dev.
ADD . /app

# Add runit service configuration
ADD docker/service /etc/service

# Listen on port 80, accept from anywhere
ENV LISTEN 0.0.0.0:80

EXPOSE 80

# Boot my_init by default; starts services, etc.
CMD ["my_init"]
