ARG BASE_IMAGE=alpine:3.15

# Final Docker image
FROM ${BASE_IMAGE} AS final-stage
LABEL MAINTAINER "Thomas Labarussias <issif+falcosidekick@gadz.org>"
RUN apk add --update --no-cache ca-certificates
# Create user falcosidekick
RUN addgroup -S falcosidekickui && adduser -u 1234 -S falcosidekickui -G falcosidekickui
# must be numeric to work with Pod Security Policies:
# https://kubernetes.io/docs/concepts/policy/pod-security-policy/#users-and-groups
USER 1234
WORKDIR /app
COPY frontend/dist frontend/dist
COPY LICENSE .
COPY falcosidekick-ui .

EXPOSE 2802
ENTRYPOINT ["./falcosidekick-ui"]
