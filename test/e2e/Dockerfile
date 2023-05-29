FROM alpine:3.17

RUN set -x && apk add curl jq postgresql-client

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
  install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

COPY ./run.sh /test/e2e/
COPY ./manifest/*.yaml /test/e2e/manifest/

WORKDIR /test/e2e

ENTRYPOINT ["./run.sh"]
