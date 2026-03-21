# Smart Router Configuration

**Package**: `pkg/routing/` · **Main Doc**: [README.md](./README.md)

## Configuration Overview

The Smart Router supports two modes:

1. **4-Tier Mode** (V2) - Recommended for fine-grained control
2. **2-Tier Mode** (Legacy) - Backward compatible

## 4-Tier Mode Configuration

### Basic Configuration

```json
{
  "agents": {
    "defaults": {
      "routing": {
        "enabled": true,
        "tier_mapping": {
          "simple": "model-for-simple",
          "moderate": "model-for-moderate",
          "complex": "model-for-complex",
          "reasoning": "model-for-reasoning"
        }
      }
    }
  }
}
```

### Configuration Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `enabled` | bool | No | Enable/disable routing (default: false) |
| `tier_mapping` | map | Yes* | Maps tier names to model names |
| `tier_boundaries` | object | No | Custom score boundaries for tiers |

*Either `tier_mapping` (4-tier) or `light_model` (2-tier) is required when enabled.

### Tier Mapping

The `tier_mapping` field maps each tier to a model name from your `model_list`:

```json
{
  "tier_mapping": {
    "simple": "glm-4-flash",
    "moderate": "glm-4-flash",
    "complex": "glm-4-plus",
    "reasoning": "claude-opus-4-6"
  }
}
```

**Valid tier keys**: `simple`, `moderate`, `complex`, `reasoning` (case-sensitive)

**Invalid keys**: Logged and ignored

### Custom Tier Boundaries

Optionally customize the score boundaries for each tier:

```json
{
  "routing": {
    "enabled": true,
    "tier_mapping": {
      "simple": "model-a",
      "moderate": "model-b",
      "complex": "model-c",
      "reasoning": "model-d"
    },
    "tier_boundaries": {
      "simple_moderate": 0.0,
      "moderate_complex": 0.25,
      "complex_reasoning": 0.5
    }
  }
}
```

#### Boundary Validation Rules

1. All cutpoints must be within `[-1.0, 1.0]`
2. Must satisfy: `simple_moderate < moderate_complex < complex_reasoning`
3. Invalid configurations fall back to defaults and log a warning

#### Default Boundaries

| Cutpoint | Default Value |
|----------|---------------|
| `simple_moderate` | -0.05 |
| `moderate_complex` | 0.15 |
| `complex_reasoning` | 0.35 |

## 2-Tier Mode Configuration (Legacy)

For backward compatibility, the original 2-tier mode is still supported:

```json
{
  "agents": {
    "defaults": {
      "routing": {
        "enabled": true,
        "light_model": "glm-4-flash",
        "threshold": 0.35
      }
    }
  }
}
```

### 2-Tier Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `light_model` | string | - | Model for simple tasks |
| `threshold` | float | 0.35 | Score cutoff (≥ threshold → heavy model) |

## Mode Detection Logic

The router automatically detects which mode to use:

```
if tier_mapping is configured:
    → Use 4-tier mode (V2)
elif light_model is configured:
    → Use 2-tier mode (legacy)
else:
    → Disable routing, use primary model
```

## Configuration Examples

### Cost-Optimized Setup

Route simple tasks to the cheapest models:

```json
{
  "routing": {
    "enabled": true,
    "tier_mapping": {
      "simple": "gpt-3.5-turbo",
      "moderate": "gpt-3.5-turbo",
      "complex": "gpt-4",
      "reasoning": "gpt-4"
    }
  }
}
```

### Performance-Optimized Setup

Use fast models for most tasks, powerful models for complex work:

```json
{
  "routing": {
    "enabled": true,
    "tier_mapping": {
      "simple": "groq/llama-3.1-70b",
      "moderate": "groq/llama-3.1-70b",
      "complex": "claude-sonnet-4-6",
      "reasoning": "claude-opus-4-6"
    }
  }
}
```

### Multi-Provider Setup

Mix providers based on task type:

```json
{
  "routing": {
    "enabled": true,
    "tier_mapping": {
      "simple": "zhipu/glm-4-flash",
      "moderate": "zhipu/glm-4-flash",
      "complex": "openai/gpt-4",
      "reasoning": "anthropic/claude-opus-4-6"
    }
  }
}
```

### Custom Boundaries for Aggressive Routing

Route more messages to simpler models:

```json
{
  "routing": {
    "enabled": true,
    "tier_mapping": {
      "simple": "model-a",
      "moderate": "model-b",
      "complex": "model-c",
      "reasoning": "model-d"
    },
    "tier_boundaries": {
      "simple_moderate": 0.0,
      "moderate_complex": 0.3,
      "complex_reasoning": 0.6
    }
  }
}
```

## Environment Variables

The router respects these environment variables:

| Variable | Description |
|----------|-------------|
| `MOONHUB_ROUTING_ENABLED` | Override routing enabled setting |
| `MOONHUB_ROUTING_LIGHT_MODEL` | Override light model |
| `MOONHUB_ROUTING_THRESHOLD` | Override threshold |

## Configuration Migration

To migrate from 2-tier to 4-tier:

```json
// Before (2-tier)
{
  "routing": {
    "enabled": true,
    "light_model": "glm-4-flash",
    "threshold": 0.35
  }
}

// After (4-tier)
{
  "routing": {
    "enabled": true,
    "tier_mapping": {
      "simple": "glm-4-flash",
      "moderate": "glm-4-flash",
      "complex": "glm-4-plus",
      "reasoning": "claude-opus-4-6"
    }
  }
}
```

## Validation Errors

| Error | Cause | Resolution |
|-------|-------|------------|
| "unknown tier key" | Invalid tier name in mapping | Use: simple, moderate, complex, reasoning |
| "invalid boundaries" | Cutpoints out of order or range | Ensure sm < mc < cr within [-1, 1] |
| "no model configured" | Empty tier mapping | Add model names to tier_mapping |
