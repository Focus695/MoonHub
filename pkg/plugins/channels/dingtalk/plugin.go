// MoonHub - Ultra-lightweight personal AI agent
// Plugin Architecture - DingTalk Channel Plugin
//
// Copyright (c) 2026 MoonHub contributors

package dingtalk

import (
	"fmt"

	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	channelsdingtalk "github.com/sipeed/moonhub/pkg/channels/dingtalk"
	"github.com/sipeed/moonhub/pkg/config"
	"github.com/sipeed/moonhub/pkg/framework"
)

func init() {
    plugin.RegisterPlugin(&DingTalkPlugin{})
}

// DingTalkPlugin implements ChannelPlugin interface for DingTalk
type DingTalkPlugin struct{}

// Metadata returns plugin metadata
func (p *DingTalkPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        ID:          "moonhub-channel-dingtalk",
        Name:        "DingTalk",
        Type:        plugin.TypeChannel,
        Version:     "2.0.0",
        Description: "DingTalk bot channel integration",
        Priority:    100,
    }
}

// Init initializes the plugin with runtime context
func (p *DingTalkPlugin) Init(ctx *plugin.RuntimeContext) error {
    return nil
}

// Validate checks if the plugin can run with current config
func (p *DingTalkPlugin) Validate(cfg *config.Config) error {
    if cfg.Channels.DingTalk.Enabled && cfg.Channels.DingTalk.ClientID == "" {
        return fmt.Errorf("dingtalk client_id required when enabled")
    }
    return nil
}

// ChannelPrefix returns the userId prefix this channel owns
func (p *DingTalkPlugin) ChannelPrefix() string {
    return "dingtalk"
}

// IsEnabled checks if the channel is enabled in config
func (p *DingTalkPlugin) IsEnabled(cfg *config.Config) bool {
    return cfg.Channels.DingTalk.Enabled && cfg.Channels.DingTalk.ClientID != ""
}
// CreateChannel instantiates the channel implementation
func (p *DingTalkPlugin) CreateChannel(cfg *config.Config, bus *bus.MessageBus) (channels.Channel, error) {
	return channelsdingtalk.NewDingTalkChannel(cfg.Channels.DingTalk, bus)
}
