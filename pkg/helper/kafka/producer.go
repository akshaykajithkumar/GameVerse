// kafka.go
package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

var (
	producer sarama.SyncProducer
	consumer sarama.Consumer
	admin    sarama.ClusterAdmin
)

func init() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	var err error

	producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err = sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}

	admin, err = sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
}

func createTopicIfNotExists(topic string) error {
	// Check if the topic exists
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	if _, exists := topics[topic]; exists {
		// Topic already exists
		return nil
	}

	// Topic doesn't exist, create it
	config := sarama.NewConfig()
	config.Admin.Retry.Max = 5

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		return err
	}

	// Wait for the topic to be created
	for i := 0; i < 30; i++ {
		topics, err := admin.ListTopics()
		if err != nil {
			return err
		}

		if _, exists := topics[topic]; exists {
			// Topic created successfully
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for topic creation")
}

func ProduceToKafka(userID, videoID int, videoURL string) error {
	topic := fmt.Sprintf("%s_%d_%d", "user_topic", userID, videoID)

	// Check if the topic exists, create it if not
	if err := createTopicIfNotExists(topic); err != nil {
		return fmt.Errorf("error creating topic: %v", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(videoURL),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		return err
	}

	fmt.Printf("Produced message to topic %s (partition %d, offset %d)\n", topic, partition, offset)
	return nil
}
