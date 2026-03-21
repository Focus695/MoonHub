# MoonHub Learning System Configuration

This document describes the configuration options for the MoonHub learning system.

## Overview

The learning system supports various configuration options to customize behavior.

## Configuration Structure

Configure the learning system in `~/.moonhub/config.json`:
```json
{
  "learning": {
    "language": "auto",
    "db_path": "memory/learning.db",
    "min_confidence": 0.7,
    "max_examples": 10,
    "enable_semantic_detection": true,
    "enable_implicit_signals": true,
    "enable_behavioral_scoring": true,
    "enable_pattern_evolution": false,
    "enable_contradiction_detection": false,
    "enable_suggestions": false,
    "decay_older_than_days": 7,
    "prune_older_than_days": 30,
    "merge_similarity_threshold": 0.8
  }
}
```

## Configuration Options

### language
- `auto` - Auto-detect language (default)
- `en` - Use English mode only
- `zh` - Use Chinese mode only

### db_path
Learning database storage path, defaults to `memory/learning.db`

### min_confidence
Pattern confidence threshold; patterns below this value will be ignored, default 0.7

### max_examples
Maximum number of examples to retain per pattern, default 10

### enable_semantic_detection
Enable semantic keyword detection, default true

### enable_implicit_signals
Infer implicit signals from tool usage, default true

### enable_behavioral_scoring
Enable behavioral scoring system, default true

### enable_pattern_evolution
Enable pattern evolution (merge/prune/decay), default false

### enable_contradiction_detection
Detect contradictory patterns, default false

### enable_suggestions
Generate proactive suggestions, default false

### decay_older_than_days
Pattern decay period in days, default 7 days

### prune_older_than_days
Pattern pruning period in days, default 30 days

### merge_similarity_threshold
Pattern merge similarity threshold, default 0.8

## Usage Examples

### Chinese Mode Configuration
```json
{
  "learning": {
    "language": "zh"
  }
}
```

### English Mode Configuration
```json
{
  "learning": {
    "language": "en"
  }
}
```

### Customized Configuration
```json
{
  "learning": {
    "language": "auto",
    "enable_suggestions": true,
    "enable_pattern_evolution": true
  }
}
```

## Notes
- Configuration changes require service restart to take effect
- `db_path` directory requires write permissions
- Advanced features (evolution, contradiction detection) may affect performance
