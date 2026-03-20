// MoonHub - Ultra-lightweight personal AI agent
// Plugin Architecture - WhatsApp Channel Plugin
//
// Copyright (c) 2026 MoonHub contributors

package whatsapp

import (
	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/channels"
	channelswhatsapp "github.com/sipeed/moonhub/pkg/channels/whatsapp"
	"github.com/sipeed/moonhub/pkg/config"
	"github.com/sipeed/moonhub/pkg/framework"
)

func init() {
    plugin.RegisterPlugin(&WhatsAppPlugin{})
}

// WhatsAppPlugin implements ChannelPlugin interface for WhatsApp
type WhatsAppPlugin struct{}

// Metadata returns plugin metadata
func (p *WhatsAppPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        ID:          "moonhub-channel-whatsapp",
        Name:        "WhatsApp",
        Type:        plugin.TypeChannel,
        Version:     "2.0.0",
        Description: "WhatsApp channel integration via bridge",
        Priority:    100,
    }
}

// Init initializes the plugin with runtime context
func (p *WhatsAppPlugin) Init(ctx *plugin.RuntimeContext) error {
    return nil
}

// Validate checks if the plugin can run with current config
func (p *WhatsAppPlugin) Validate(cfg *config.Config) error {
    return nil
}

// ChannelPrefix returns the userId prefix this channel owns
func (p *WhatsAppPlugin) ChannelPrefix() string {
    return "whatsapp"
}

// IsEnabled checks if the channel is enabled in config
func (p *WhatsAppPlugin) IsEnabled(cfg *config.Config) bool {
    wa := cfg.Channels.WhatsApp
    return wa.Enabled && !wa.UseNative && wa.BridgeURL != ""
}

// CreateChannel instantiates the channel implementation
func (p *WhatsAppPlugin) CreateChannel(cfg *config.Config, bus *bus.MessageBus) (channels.Channel, error) {
	return channelswhatsapp.NewWhatsAppChannel(cfg.Channels.WhatsApp, bus)
}
