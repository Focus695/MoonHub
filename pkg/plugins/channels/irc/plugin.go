// MoonHub - Ultra-lightweight personal AI agent
// Plugin Architecture - IRC Channel Plugin
//
// Copyright (c) 2026 MoonHub contributors

package irc

import (
	"fmt"

	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	channelsirc "github.com/sipeed/moonhub/pkg/channels/irc"
	"github.com/sipeed/moonhub/pkg/config"
	"github.com/sipeed/moonhub/pkg/framework"
)

func init() {
	plugin.RegisterPlugin(&IRCPlugin{})
}

// IRCPlugin implements ChannelPlugin interface for IRC
type IRCPlugin struct{}

// Metadata returns plugin metadata
func (p *IRCPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		ID:          "moonhub-channel-irc",
		Name:        "IRC",
		Type:        plugin.TypeChannel,
		Version:     "2.0.0",
		Description: "IRC protocol channel integration",
		Priority:    100,
	}
}

// Init initializes the plugin with runtime context
func (p *IRCPlugin) Init(ctx *plugin.RuntimeContext) error {
	return nil
}

// Validate checks if the plugin can run with current config
func (p *IRCPlugin) Validate(cfg *config.Config) error {
	if cfg.Channels.IRC.Enabled && cfg.Channels.IRC.Server == "" {
		return fmt.Errorf("irc server required when enabled")
	}
	return nil
}

// ChannelPrefix returns the userId prefix this channel owns
func (p *IRCPlugin) ChannelPrefix() string {
	return "irc"
}

// IsEnabled checks if the channel is enabled in config
func (p *IRCPlugin) IsEnabled(cfg *config.Config) bool {
	return cfg.Channels.IRC.Enabled && cfg.Channels.IRC.Server != ""
}

// CreateChannel instantiates the channel implementation
// This calls the existing implementation from pkg/channels/irc
func (p *IRCPlugin) CreateChannel(cfg *config.Config, bus *bus.MessageBus) (channels.Channel, error) {
	return channelsirc.NewIRCChannel(cfg.Channels.IRC, bus)
}
