package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestKafkaProducer(t *testing.T) {
	// Cấu hình Kafka Producer
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "test-producer",
		"acks":              "all",
	}

	// Tạo Producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Topic để gửi tin nhắn
	topic := "test-topic"

	// Hàm xử lý kết quả gửi tin nhắn
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	// Gửi tin nhắn thử nghiệm
	message := fmt.Sprintf("Test message at %s", time.Now().String())
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		t.Fatalf("Failed to produce message: %v", err)
	}

	// Chờ phản hồi từ Kafka
	e := <-deliveryChan
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			t.Errorf("Delivery failed: %v", ev.TopicPartition.Error)
		} else {
			t.Logf("Message delivered to %v", ev.TopicPartition)
		}
	case kafka.Error:
		t.Errorf("Error: %v", ev)
	}
}

func main() {
	testing.Main(func(pat, str string) (bool, error) { return true, nil }, []testing.InternalTest{
		{Name: "TestKafkaProducer", F: TestKafkaProducer},
	}, nil, nil)
}
