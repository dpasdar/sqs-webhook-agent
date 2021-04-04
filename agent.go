package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Payload struct {
	Endpoint string              `json:"end_point"`
	Body     string              `json:"body"`
	Headers  map[string][]string `json:"headers"`
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug level messages")
	queueName := flag.String("queue_name", "", "Name/Address of the sqs queue to read from (Required).")
	webhookUrl := flag.String("webhook_url", "http://localhost:9000/hooks", "URL of the webhook(Required).")
	flag.Parse()
	if *queueName == "" {
		log.Fatal("Queue Name is required.")
		return
	}

	if *webhookUrl == "" {
		log.Fatal("Webhook URL is required.")
		return
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Infoln("Agent started...")
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		s := <-sigs
		log.Infof("RECEIVED SIGNAL: %s", s)
		os.Exit(0)
	}()

	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queue, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: queueName,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	queueURL := queue.QueueUrl
	for {
		time.Sleep(5 * time.Second)
		msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            queueURL,
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(20),
		})
		if err != nil {
			log.Error(err)
			continue
		}
		for _, message := range msgResult.Messages {
			var parsed Payload
			err := json.Unmarshal([]byte(*message.Body), &parsed)
			log.Infof("Got message with endpoint %s to be passed to webhook", parsed.Endpoint)
			if err != nil {
				log.Warn(err)
				continue
			}
			payload := parsed.Body
			endpoint := parsed.Endpoint
			headers := parsed.Headers
			req, _ := http.NewRequest("POST", *webhookUrl+"/"+endpoint, bytes.NewBufferString(payload))
			for k, v := range headers {
				log.Debugf("Header {%s: %s}", k, v)
				req.Header.Set(k, v[0])
			}
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Warn(err)
				continue
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Warn(err)
				continue
			}
			b, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				log.Error(err.Error())
				continue
			}
			log.Infof("Response from webhook: %s", string(b))
		}
	}
}
