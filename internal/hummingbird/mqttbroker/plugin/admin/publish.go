package admin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker"
	"github.com/winc-link/hummingbird/internal/pkg/packets"
)

type publisher struct {
	a *Admin
}

func (p *publisher) mustEmbedUnimplementedPublishServiceServer() {
	return
}

// Publish publishes a message into broker.
func (p *publisher) Publish(ctx context.Context, req *PublishRequest) (resp *empty.Empty, err error) {
	if !packets.ValidV5Topic([]byte(req.TopicName)) {
		return nil, ErrInvalidArgument("topic_name", "")
	}
	if req.Qos > uint32(packets.Qos2) {
		return nil, ErrInvalidArgument("qos", "")
	}
	if req.PayloadFormat != 0 && req.PayloadFormat != 1 {
		return nil, ErrInvalidArgument("payload_format", "")
	}
	if req.ResponseTopic != "" && !packets.ValidV5Topic([]byte(req.ResponseTopic)) {
		return nil, ErrInvalidArgument("response_topic", "")
	}
	var userPpt []packets.UserProperty
	for _, v := range req.UserProperties {
		userPpt = append(userPpt, packets.UserProperty{
			K: v.K,
			V: v.V,
		})
	}

	p.a.publisher.Publish(&mqttbroker.Message{
		Dup:             false,
		QoS:             byte(req.Qos),
		Retained:        req.Retained,
		Topic:           req.TopicName,
		Payload:         []byte(req.Payload),
		ContentType:     req.ContentType,
		CorrelationData: []byte(req.CorrelationData),
		MessageExpiry:   req.MessageExpiry,
		PayloadFormat:   packets.PayloadFormat(req.PayloadFormat),
		ResponseTopic:   req.ResponseTopic,
		UserProperties:  userPpt,
	})
	return &empty.Empty{}, nil
}
