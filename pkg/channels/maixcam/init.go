package maixcam

import (
	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	"github.com/sipeed/moonhub/pkg/config"
)

func init() {
	channels.RegisterFactory("maixcam", func(cfg *config.Config, b *bus.MessageBus) (channels.Channel, error) {
		return NewMaixCamChannel(cfg.Channels.MaixCam, b)
	})
}
