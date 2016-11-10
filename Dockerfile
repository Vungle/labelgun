FROM alpine

ARG VCS_REF
ARG BUILD_DATE

# Metadata
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/dailyhotel/labelgun" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.docker.dockerfile="/Dockerfile"

# Install kubectl
ENV KUBE_LATEST_VERSION="v1.4.4"

RUN apk add --update ca-certificates curl \
 && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
 && chmod +x /usr/local/bin/kubectl \
 && rm /var/cache/apk/*

# Add labelgun.go
ENV PATH /usr/local/bin:$PATH

COPY bin/labelgun /usr/local/bin/labelgun
COPY run-labelgun.sh /usr/local/bin/run-labelgun.sh

CMD ["run-labelgun.sh"]
