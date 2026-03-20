// MoonHub - Ultra-lightweight personal AI agent
// Plugin Architecture - Message Tool Plugin
//
// Copyright (c) 2026 MoonHub contributors

package message

import (
	"context"
	"time"

	"github.com/sipeed/moonhub/pkg/bus"
	"github.com/sipeed/moonhub/pkg/config"
	"github.com/sipeed/moonhub/pkg/framework"
	"github.com/sipeed/moonhub/pkg/tools"
)

func init() {
	plugin.RegisterPlugin(&MessageToolPlugin{})
}

type MessageToolPlugin struct {
	ctx *plugin.RuntimeContext
}

func (p *MessageToolPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		ID:          "moonhub-tool-message",
		Name:        "Message",
		Type:        plugin.TypeTool,
		Version:     "1.0.0",
		Description: "Send messages to chat channels",
		Priority:    100,
	}
}

func (p *MessageToolPlugin) Init(ctx *plugin.RuntimeContext) error {
	p.ctx = ctx
	return nil
}

func (p *MessageToolPlugin) Validate(cfg *config.Config) error {
	return nil
}

func (p *MessageToolPlugin) CreateTools(ctx *plugin.RuntimeContext) []tools.Tool {
	if !ctx.Config.Tools.IsToolEnabled("message") {
		return nil
	}
	msgBus := ctx.Bus
	if msgBus == nil {
		return nil
	}
	messageTool := tools.NewMessageTool()
	messageTool.SetSendCallback(func(channel, chatID, content string) error {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer pubCancel()
		return msgBus.PublishOutbound(pubCtx, bus.OutboundMessage{
			Channel: channel,
			ChatID:  chatID,
			Content: content,
		})
	})
	return []tools.Tool{messageTool}
}

func (p *MessageToolPlugin) IsCore() bool {
	return true
}
