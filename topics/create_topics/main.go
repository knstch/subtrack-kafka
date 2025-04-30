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
	Description string `yaml:"description"`
}

type TopicGroups struct {
	SubtrackTopics []Topic `yaml:"subtrack-topics"`
	DexTopics      []Topic `yaml:"dex-topics"`
}

type TopicsYAML struct {
	Topics TopicGroups `yaml:"topics"`
}

func main() {
	defaultBroker := os.Getenv("KAFKA_BROKER")
	if defaultBroker == "" {
		fmt.Println("❌ KAFKA_BROKER not set")
		os.Exit(1)
	}
	dexBroker := os.Getenv("KAFKA_DEX_BROKER")
	if dexBroker == "" {
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

	var cfg TopicsYAML
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	connDefaultBroker, err := kafka.Dial("tcp", defaultBroker)
	if err != nil {
		panic(err)
	}
	defer connDefaultBroker.Close()

	for _, topic := range cfg.Topics.SubtrackTopics {
		fmt.Println("🌀 Creating subtrack topic:", topic.Name)
		err = connDefaultBroker.CreateTopics(kafka.TopicConfig{
			Topic:             topic.Name,
			NumPartitions:     topic.Partitions,
			ReplicationFactor: topic.Replication,
		})
		if err != nil {
			fmt.Println("❌", err)
		}
	}

	err = connDefaultBroker.CreateTopics(kafka.TopicConfig{
		Topic:             "__consumer_offsets",
		NumPartitions:     50,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("⚠️ (Optional) Couldn't create __consumer_offsets: %v\n", err)
	}

	connDexBroker, err := kafka.Dial("tcp", dexBroker)
	if err != nil {
		panic(err)
	}
	defer connDexBroker.Close()

	for _, topic := range cfg.Topics.DexTopics {
		fmt.Println("🌀 Creating dex topic:", topic.Name)
		err = connDexBroker.CreateTopics(kafka.TopicConfig{
			Topic:             topic.Name,
			NumPartitions:     topic.Partitions,
			ReplicationFactor: topic.Replication,
		})
		if err != nil {
			fmt.Println("❌", err)
		}
	}

	err = connDexBroker.CreateTopics(kafka.TopicConfig{
		Topic:             "__consumer_offsets",
		NumPartitions:     50,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("⚠️ (Optional) Couldn't create __consumer_offsets: %v\n", err)
	}

	fmt.Println("✅ Done")
}
