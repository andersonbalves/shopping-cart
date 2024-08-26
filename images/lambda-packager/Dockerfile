FROM amazonlinux:latest

WORKDIR /app
RUN yum install -y golang zip
COPY apps /app/lambdas
COPY scripts/golang /app/scripts
ENTRYPOINT ["/app/scripts/package.sh"]
