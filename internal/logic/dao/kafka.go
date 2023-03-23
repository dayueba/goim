package dao

import (
	"context"
	"strconv"

	"github.com/Shopify/sarama"
	pb "github.com/Terry-Mao/goim/api/logic"
	"github.com/Terry-Mao/goim/internal/logic/conf"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type kafkaDao struct {
	kafkaPub sarama.SyncProducer
}

func (k *kafkaDao) SendMessage(topic, ackInbox string, key string, msg []byte) error {
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}
	_, _, err := k.kafkaPub.SendMessage(m)
	return err
}

func (k *kafkaDao) Close() error {
	return nil
}

func newKafka(c *conf.Kafka) PushMsg {
	return &kafkaDao{kafkaPub: newKafkaPub(c)}
}

func newKafkaPub(c *conf.Kafka) sarama.SyncProducer {
	kc := sarama.NewConfig()
	kc.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := sarama.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

// PushMsg push a message to databus.
func (d *Dao) PushMsg(c context.Context, op int32, server string, keys []string, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Type:      pb.PushMsg_PUSH,
		Operation: op,
		Server:    server,
		Keys:      keys,
		Msg:       msg,
	}

	// todo
	// 即时消息存储扩展 HOOKS:
	// 在这里增加即时消息存储扩展
	// 如果需要只存储离线消息, 可以先检查当前用户是否在线, 依据用户在线情况处理存储

	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	if err = d.mq.SendMessage(d.c.Kafka.Topic, "", keys[0], b); err != nil {
		log.Errorf("PushMsg.send(push pushMsg:%v) error(%v)", pushMsg, err)
	}
	return
}

// BroadcastRoomMsg push a message to databus.
func (d *Dao) BroadcastRoomMsg(c context.Context, op int32, room string, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Type:      pb.PushMsg_ROOM,
		Operation: op,
		Room:      room,
		Msg:       msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	if err = d.mq.SendMessage(d.c.Kafka.Topic, "", room, b); err != nil {
		log.Errorf("PushMsg.send(broadcast_room pushMsg:%v) error(%v)", pushMsg, err)
	}
	return
}

// BroadcastMsg push a message to databus.
func (d *Dao) BroadcastMsg(c context.Context, op, speed int32, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Type:      pb.PushMsg_BROADCAST,
		Operation: op,
		Speed:     speed,
		Msg:       msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	if err = d.mq.SendMessage(d.c.Kafka.Topic, "", strconv.FormatInt(int64(op), 10), b); err != nil {
		log.Errorf("PushMsg.send(broadcast pushMsg:%v) error(%v)", pushMsg, err)
	}
	return
}
