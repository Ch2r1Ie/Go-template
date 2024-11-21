package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Example:
// defer func() {
// 	if err := producer.Close(); err != nil {
// 		log.Panic(err)
// 	}
// }()

// msg := &sarama.ProducerMessage{Topic: "my_topic", Value: sarama.StringEncoder("testing 123")}
// partition, offset, err := producer.SendMessage(msg)
// if err != nil {
// 	log.Printf("FAILED to send message: %s\n", err)
// } else {
// 	log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
// }

func NewSyncProducerGuarantee(addrs []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		log.Panic(err)
	}
	return producer
}

func NewSyncProducerFirenForget(addrs []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		log.Panic(err)
	}
	return producer
}
