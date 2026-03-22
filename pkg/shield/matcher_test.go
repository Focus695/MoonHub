package shield

import (
	"testing"
)

func TestMatchToolCall(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		name        string
		condition   string
		event       ShieldEvent
		shouldMatch bool
	}{
		{
			name:        "exact tool name match",
			condition:   "exec",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolName: "exec"},
			shouldMatch: true,
		},
		{
			name:        "partial tool name match",
			condition:   "exec",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolName: "execute_code"},
			shouldMatch: true,
		},
		{
			name:        "no tool name match",
			condition:   "exec",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolName: "read_file"},
			shouldMatch: false,
		},
		{
			name:        "case insensitive match",
			condition:   "EXEC",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolName: "exec"},
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.matchToolCall(tt.condition, tt.event)
			matched := result != nil
			if matched != tt.shouldMatch {
				t.Errorf("matchToolCall() matched = %v, expected %v", matched, tt.shouldMatch)
			}
		})
	}
}

func TestMatchToolArgsContaining(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		name        string
		patterns    string
		event       ShieldEvent
		shouldMatch bool
	}{
		{
			name:        "SQL injection pattern match",
			patterns:    "(DROP, DELETE, UNION)",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolArgs: map[string]any{"command": "DROP TABLE users"}},
			shouldMatch: true,
		},
		{
			name:        "no pattern match",
			patterns:    "(DROP, DELETE)",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolArgs: map[string]any{"command": "SELECT * FROM users"}},
			shouldMatch: false,
		},
		{
			name:        "empty args",
			patterns:    "(DROP, DELETE)",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolArgs: map[string]any{}},
			shouldMatch: false,
		},
		{
			name:        "case insensitive match",
			patterns:    "(drop, delete)",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolArgs: map[string]any{"query": "DROP TABLE"}},
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.matchToolArgsContaining(tt.patterns, tt.event)
			matched := result != nil
			if matched != tt.shouldMatch {
				t.Errorf("matchToolArgsContaining() matched = %v, expected %v", matched, tt.shouldMatch)
			}
		})
	}
}

func TestMatchFilePathContains(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		name        string
		path        string
		event       ShieldEvent
		shouldMatch bool
	}{
		{
			name:        "path traversal match",
			path:        "../",
			event:       ShieldEvent{Scope: ScopeToolCall, FilePath: "../../../etc/passwd"},
			shouldMatch: true,
		},
		{
			name:        "sensitive file match",
			path:        "/etc/passwd",
			event:       ShieldEvent{Scope: ScopeToolCall, FilePath: "/etc/passwd"},
			shouldMatch: true,
		},
		{
			name:        "no match",
			path:        "/etc/passwd",
			event:       ShieldEvent{Scope: ScopeToolCall, FilePath: "/home/user/file.txt"},
			shouldMatch: false,
		},
		{
			name:        "match from tool args",
			path:        ".env",
			event:       ShieldEvent{Scope: ScopeToolCall, ToolArgs: map[string]any{"path": "/app/.env"}},
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.matchFilePathContains(tt.path, tt.event)
			matched := result != nil
			if matched != tt.shouldMatch {
				t.Errorf("matchFilePathContains() matched = %v, expected %v", matched, tt.shouldMatch)
			}
		})
	}
}

func TestMatchEvent(t *testing.T) {
	matcher := NewMatcher()

	threats := []ThreatEntry{
		{
			ID:                  "THREAT-001",
			Category:            CategoryTool,
			Severity:            SeverityHigh,
			Confidence:          0.9,
			Action:              ActionBlock,
			RecommendationAgent: "BLOCK: tool.call exec\nBLOCK: tool.call with arguments containing (DROP, DELETE)",
		},
		{
			ID:                  "THREAT-002",
			Category:            CategoryFile,
			Severity:            SeverityHigh,
			Confidence:          0.9,
			Action:              ActionBlock,
			RecommendationAgent: "BLOCK: file path contains /etc/passwd\nBLOCK: file path contains ../",
		},
	}

	tests := []struct {
		name       string
		event      ShieldEvent
		matchCount int
	}{
		{
			name:       "match tool name",
			event:      ShieldEvent{Scope: ScopeToolCall, ToolName: "exec"},
			matchCount: 1,
		},
		{
			name:       "match tool args",
			event:      ShieldEvent{Scope: ScopeToolCall, ToolName: "db_query", ToolArgs: map[string]any{"query": "DROP TABLE"}},
			matchCount: 1,
		},
		{
			name:       "match file path",
			event:      ShieldEvent{Scope: ScopeToolCall, FilePath: "/etc/passwd"},
			matchCount: 1,
		},
		{
			name:       "no match",
			event:      ShieldEvent{Scope: ScopeToolCall, ToolName: "read_file", ToolArgs: map[string]any{"path": "/home/user/file.txt"}},
			matchCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := matcher.MatchEvent(tt.event, threats)
			if len(results) != tt.matchCount {
				t.Errorf("MatchEvent() returned %d matches, expected %d", len(results), tt.matchCount)
			}
		})
	}
}

func TestIsScopeCompatible(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		scope    ShieldScope
		category ThreatCategory
		expected bool
	}{
		{ScopeToolCall, CategoryTool, true},
		{ScopeToolCall, CategoryFile, true},
		{ScopeToolCall, CategoryNetwork, false},
		{ScopeNetworkEgress, CategoryNetwork, true},
		{ScopeNetworkEgress, CategoryTool, false},
		{ScopeSkillExecute, CategorySkill, true},
		{ScopePrompt, CategoryPrompt, true},
		{ScopePrompt, CategoryMessage, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.scope)+"_"+string(tt.category), func(t *testing.T) {
			result := matcher.isScopeCompatible(tt.scope, tt.category)
			if result != tt.expected {
				t.Errorf("isScopeCompatible(%s, %s) = %v, expected %v", tt.scope, tt.category, result, tt.expected)
			}
		})
	}
}

func TestExtractPatternsFromParens(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "(DROP, DELETE, UNION)",
			expected: []string{"DROP", "DELETE", "UNION"},
		},
		{
			input:    "( single )",
			expected: []string{"single"},
		},
		{
			input:    "no parens",
			expected: []string{"no parens"},
		},
		{
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractPatternsFromParens(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("extractPatternsFromParens() = %v, expected %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("extractPatternsFromParens()[%d] = %q, expected %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestSplitByOr(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "domain1.com or domain2.com",
			expected: []string{"domain1.com", "domain2.com"},
		},
		{
			input:    "a OR b Or c",
			expected: []string{"a", "b", "c"},
		},
		{
			input:    "single-domain.com",
			expected: []string{"single-domain.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := splitByOr(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitByOr() = %v, expected %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("splitByOr()[%d] = %q, expected %q", i, v, tt.expected[i])
				}
			}
		})
	}
}
