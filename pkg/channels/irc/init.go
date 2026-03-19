package irc

import (
	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	"github.com/sipeed/moonhub/pkg/config"
)

func init() {
	channels.RegisterFactory("irc", func(cfg *config.Config, b *bus.MessageBus) (channels.Channel, error) {
		if !cfg.Channels.IRC.Enabled {
			return nil, nil
		}
		return NewIRCChannel(cfg.Channels.IRC, b)
	})
}
