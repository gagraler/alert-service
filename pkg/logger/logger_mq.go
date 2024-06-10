package logger

import (
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap/zapcore"
	"strings"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/8 1:13
 * @file: logger_mq.go
 * @description: log write mq
 */

type kafkaWriter struct {
	producer sarama.SyncProducer
	topic    string
}

func (k *kafkaWriter) Write(p []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(p),
	}

	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (k *kafkaWriter) Sync() error {
	return nil
}

// getKafkaLogWriter returns the WriteSyncer for logging to Kafka.
func getKafkaLogWriter(config *Config) (zapcore.WriteSyncer, error) {
	producer, err := sarama.NewSyncProducer(strings.Split(config.KafkaBrokers, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	return zapcore.AddSync(&kafkaWriter{
		producer: producer,
		topic:    config.KafkaTopic,
	}), nil
}
