package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"rmq"
	"dao"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULT_EXCHANGE      string = "popupflow_exchange"
	DEFAULT_EXCHANGE_TYPE string = "topic"
	DEFAULT_QUEUE         string = "popupflow_queue"
	DEFAULT_ROUTING_KEY   string = "popupflow.activity.1111"
	DEFAULT_URI           string = "amqp://guest:guest@localhost:5672/"
	DEFAULT_CONCURRENCY   int    = 5
)

// Flags
var (
	consumer = flag.Bool("c", true, "Act as a consumer")
	producer = flag.Bool("p", false, "Act as a producer")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  Consumer : amqpc -c\n")
	fmt.Fprintf(os.Stderr, "  Producer : amqpc -p\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func main() {

	done := make(chan error)

	flag.Usage = usage

	if *producer {
		test := dao.RechargeInfo{
			CompereId: 123,
			UserId: 234,
			Sid: 345,
			Ssid: 456,
			SignSid: 567,
			GiftId: 678,
			GiftValue: 789,
			Platform: 8910,
			Timestamp: 91011,
			Days: "2016:10:17",
			Hours: "15:00",
			Group: "YY",
		}
		body, err := bson.Marshal(test)
		if err != nil {
			log.Println("failed to marshal to bson")
			return
		}
		for i := 0; i < DEFAULT_CONCURRENCY; i++ {
			go startProducer(done, body)
		}

	} else if *consumer {
		for i := 0; i < DEFAULT_CONCURRENCY; i++ {
			go startConsumer(done)
		}
	}

	err := <-done
	if err != nil {
		log.Fatalf("Error : %s", err)
	}

	log.Printf("Exiting...")
}

func startConsumer(done chan error) {

	_, err := rmq.NewConsumer(
		DEFAULT_URI,
		DEFAULT_EXCHANGE,
		DEFAULT_EXCHANGE_TYPE,
		DEFAULT_QUEUE,
		DEFAULT_ROUTING_KEY,
	)

	if err != nil {
		log.Fatalf("Error while starting consumer : %s", err)
	}

	<-done
}

func startProducer(done chan error, body []byte) {
	p, err := rmq.NewProducer(
		DEFAULT_URI,
		DEFAULT_EXCHANGE,
		DEFAULT_EXCHANGE_TYPE,
	)

	if err != nil {
		log.Fatalf("Error while starting producer : %s", err)
	}

	for {
		p.Publish(DEFAULT_EXCHANGE, DEFAULT_ROUTING_KEY, body)
	}

	done <- nil
}

