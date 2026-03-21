# SHIELD Usage Examples

## Providing Custom `SHIELD.md` in Workspace

Place `SHIELD.md` in the Agent workspace root directory, for example to add a custom tool block:

````markdown
## Custom

```yaml
id: THREAT-CUSTOM-001
fingerprint: block-dangerous-exec
category: tool
severity: high
confidence: 0.9
action: block
title: Block specific tool
description: Example custom rule
recommendation_agent: |
  BLOCK: tool.call dangerous_internal_tool
```
````

If the file doesn't exist, the runtime uses the built-in default threat library (see built-in table in implementation status document).

## Programmatic Usage (Unit Tests or Standalone Tools)

```go
engine := shield.NewEngineWithDefaults()
// or shield.NewEngineFromFileOrDefault(filepath.Join(workspace, "SHIELD.md"))

decision := engine.Evaluate(shield.ShieldEvent{
    Scope:    shield.ScopeToolCall,
    ToolName: "some_tool",
    ToolArgs: map[string]any{"x": "y"},
})
// decision.Action: block / require_approval / log
```

Hot reload: `engine.Reload(content)` or `engine.ReloadFromFile(path)`.

## Runtime Approval (User Side)

When a policy matches and is downgraded or configured as `require_approval`, the session may prompt for an approval ID. Use:

- `/approve <approval_id>`
- `/reject <approval_id>`

The flow and message format follow `pkg/agent/loop.go` and the implementation status document.
