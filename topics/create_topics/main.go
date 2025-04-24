package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Topic struct {
	Name        string `yaml:"name"`
	Partitions  int    `yaml:"partitions"`
	Replication int    `yaml:"replication"`
}

type TopicFile struct {
	Topics []Topic `yaml:"topics"`
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		fmt.Println("❌ KAFKA_BROKER not set")
		os.Exit(1)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("unable to get current file path")
	}

	scriptDir := filepath.Dir(filename)
	yamlPath := filepath.Join(scriptDir, "..", "topics.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		panic(err)
	}

	var cfg TopicFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for _, topic := range cfg.Topics {
		fmt.Println("🌀 Creating topic:", topic.Name)
		err := conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic.Name,
			NumPartitions:     topic.Partitions,
			ReplicationFactor: topic.Replication,
		})
		if err != nil {
			fmt.Println("❌", err)
		}
	}

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             "__consumer_offsets",
		NumPartitions:     50,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("⚠️ (Optional) Couldn't create __consumer_offsets: %v\n", err)
	}

	fmt.Println("✅ Done")
}
