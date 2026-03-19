package feishu

import (
	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	"github.com/sipeed/moonhub/pkg/config"
)

func init() {
	channels.RegisterFactory("feishu", func(cfg *config.Config, b *bus.MessageBus) (channels.Channel, error) {
		return NewFeishuChannel(cfg.Channels.Feishu, b)
	})
}
