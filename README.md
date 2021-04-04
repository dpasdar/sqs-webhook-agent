# Webhook SQS-Relay Agent

This repository contains the complimentary solution to the [Sender](https://github.com/dpasdar/sqs-webhook-sender), which has to be installed on the every machine alongside the [webhook](https://github.com/adnanh/webhook). Please see the documentation for [Sender](https://github.com/dpasdar/sqs-webhook-sender) for more information about the architecutre.

# Installation
```
go get -u github.com/dpasdar/sqs-webhook-agent
```

In the style of golang's self-contained binaries, it will create one binary which can be copied/installed anywhere.

Optional, one can create a systemd unit file. A copy of the unit file can be found under the template/ directory. Follow the step below to create the service:

1. Copy the executable to /opt/sqs-webhook-agent/
2. Create an .env file in /opt/sqs-webhook-agent with the AWS credentials
3. Run the following to create the service (adjust the queue name in the template)
```
$ sudo useradd sqs-webhook-agent -s /sbin/nologin -M
$ sudo mkdir /opt/sqs-webhook-agent
$ sudo cp sqs/opt/sqs-webhook-agent
$ sudo cp sqs-webhook-agent.service /lib/systemd/system/.
$ sudo chmod 755 /lib/systemd/system/sqs-webhook-agent.service
$ sudo systemctl enable sqs-webhook-agent.service
```


# Running the Agent
Please make sure the following env variables are provided. Ideally, the access keys should only allow sqs operations:
```
AWS_ACCESS_KEY_ID
AWS_REGION
AWS_SECRET_ACCESS_KEY
```

Run the agent as follows:

```
$ sqs-webhook-agent -queue_name "some name"
```