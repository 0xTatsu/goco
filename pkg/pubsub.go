package pkg

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type PubSub struct {
	pubSubTopic *pubsub.Topic
}

func NewPubSub(pubSubTopic *pubsub.Topic) *PubSub {
	return &PubSub{
		pubSubTopic: pubSubTopic,
	}
}

func (r *PubSub) Publish(ctx context.Context, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &pubsub.Message{Data: dataJSON}

	publishResult := r.pubSubTopic.Publish(ctx, msg)
	if _, err = publishResult.Get(ctx); err != nil {
		return err
	}

	return nil
}
