package producer

import (
	"kafka/model"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

func SentToKafka(msg string) error {

	config := model.GetConfiguration()
	kafkaConfig := getKafkaConfig("", "")
	producers, err := sarama.NewSyncProducer([]string{config.KafkaConnection}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Unable to create kafka producer got error %v", err)
		return err
	}
	defer func() {
		if err := producers.Close(); err != nil {
			logrus.Errorf("Unable to stop kafka producer: %v", err)
			return
		}
	}()

	logrus.Infof("Success create kafka sync-producer")

	kafka := &KafkaProducer{
		Producer: producers,
	}

	if err = kafka.SendMessage("test_topic", msg); err != nil {
		logrus.Errorf("Failed send message to kafka: %v", err)
		return err
	}
	return nil
}

// KafkaProducer hold kafka producer session
type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// SendMessage function to send message into kafka
func (p *KafkaProducer) SendMessage(topic, msg string) error {

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := p.Producer.SendMessage(kafkaMsg)
	if err != nil {
		//logrus.Errorf("Send message error: %v", err)
		return err
	}

	logrus.Infof("Send message success, Topic %v, Partition %v, Offset %d", topic, partition, offset)
	return nil
}

func getKafkaConfig(username, password string) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0

	if username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	}
	return kafkaConfig
}
