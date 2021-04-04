FROM alpine
COPY sqs-webhook-agent /sqs-webhook-agent
ENTRYPOINT ["/sqs-webhook-agent"]
