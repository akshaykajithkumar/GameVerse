package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

func ConsumeURLFromKafkaForUser(ctx context.Context, userID, videoID int, urlChannel chan<- string) {
	topic := fmt.Sprintf("%s_%d_%d", "user_topic", userID, videoID)

	// Check if the topic exists, create it if not
	if err := createTopicIfNotExists(topic); err != nil {
		log.Printf("Error creating topic: %v", err)
		return
	}

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Printf("Error getting partitions: %v", err)
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	wg := &sync.WaitGroup{}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Printf("Error creating partition consumer: %v", err)
			return
		}
		defer partitionConsumer.Close()

		wg.Add(1)
		go func(pc sarama.PartitionConsumer, p int32) {
			defer wg.Done()
			for {
				select {
				case msg := <-pc.Messages():
					// Process the received message (video URL in this case)
					urlChannel <- string(msg.Value)
				case err := <-pc.Errors():
					log.Printf("Error: %v\n", err)
				case <-signals:
					return
				case <-ctx.Done():
					return
				}
			}
		}(partitionConsumer, partition)
	}

	wg.Wait()
	close(urlChannel)
}
