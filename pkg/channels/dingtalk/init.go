package dingtalk

import (
	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	"github.com/sipeed/moonhub/pkg/config"
)

func init() {
	channels.RegisterFactory("dingtalk", func(cfg *config.Config, b *bus.MessageBus) (channels.Channel, error) {
		return NewDingTalkChannel(cfg.Channels.DingTalk, b)
	})
}
