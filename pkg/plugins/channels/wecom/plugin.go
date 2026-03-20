// MoonHub - Ultra-lightweight personal AI agent
// Plugin Architecture - WeCom Channel Plugin
//
// Copyright (c) 2026 MoonHub contributors

package wecom

import (
	"fmt"

	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	channelswecom "github.com/sipeed/moonhub/pkg/channels/wecom"
	"github.com/sipeed/moonhub/pkg/config"
	"github.com/sipeed/moonhub/pkg/framework"
)

func init() {
    plugin.RegisterPlugin(&WeComPlugin{})
}

// WeComPlugin implements ChannelPlugin interface for WeCom
type WeComPlugin struct{}

// Metadata returns plugin metadata
func (p *WeComPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        ID:          "moonhub-channel-wecom",
        Name:        "WeCom",
        Type:        plugin.TypeChannel,
        Version:     "2.0.0",
        Description: "WeCom bot channel integration",
        Priority:    100,
    }
}

// Init initializes the plugin with runtime context
func (p *WeComPlugin) Init(ctx *plugin.RuntimeContext) error {
    return nil
}

// Validate checks if the plugin can run with current config
func (p *WeComPlugin) Validate(cfg *config.Config) error {
    if cfg.Channels.WeCom.Enabled && cfg.Channels.WeCom.Token == "" {
        return fmt.Errorf("wecom token required when enabled")
    }
    return nil
}

// ChannelPrefix returns the userId prefix this channel owns
func (p *WeComPlugin) ChannelPrefix() string {
    return "wecom"
}

// IsEnabled checks if the channel is enabled in config
func (p *WeComPlugin) IsEnabled(cfg *config.Config) bool {
    return cfg.Channels.WeCom.Enabled && cfg.Channels.WeCom.Token != ""
}

// CreateChannel instantiates the channel implementation
func (p *WeComPlugin) CreateChannel(cfg *config.Config, bus *bus.MessageBus) (channels.Channel, error) {
	return channelswecom.NewWeComBotChannel(cfg.Channels.WeCom, bus)
}
