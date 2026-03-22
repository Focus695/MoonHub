package shield

import (
	"testing"
	"time"
)

func TestParseYAMLBlocks(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "empty content",
			content:  "",
			expected: 0,
		},
		{
			name:     "no code blocks",
			content:  "Some text without code blocks",
			expected: 0,
		},
		{
			name: "one yaml block",
			content: `
# Some header
` + "```yaml" + `
id: THREAT-001
category: tool
severity: high
confidence: 0.9
action: block
title: Test Threat
` + "```" + `
`,
			expected: 1,
		},
		{
			name: "multiple yaml blocks",
			content: `
` + "```yaml" + `
id: THREAT-001
category: tool
severity: high
confidence: 0.9
action: block
title: Test Threat 1
` + "```" + `

` + "```yaml" + `
id: THREAT-002
category: file
severity: medium
confidence: 0.8
action: log
title: Test Threat 2
` + "```" + `
`,
			expected: 2,
		},
		{
			name: "non-threat yaml block",
			content: `
` + "```yaml" + `
name: something
value: test
` + "```" + `
`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			threats := parser.Parse(tt.content)
			if len(threats) != tt.expected {
				t.Errorf("Parse() returned %d threats, expected %d", len(threats), tt.expected)
			}
		})
	}
}

func TestParseThreatBlock(t *testing.T) {
	parser := NewParser()

	block := `id: THREAT-001
fingerprint: sql-inject-01
category: tool
severity: high
confidence: 0.90
action: block
title: SQL Injection
description: Block SQL injection attempts
recommendation_agent: |
  BLOCK: tool.call with arguments containing (DROP, DELETE)
`

	threat := parser.parseThreatBlock(block)
	if threat == nil {
		t.Fatal("parseThreatBlock returned nil")
	}

	if threat.ID != "THREAT-001" {
		t.Errorf("ID = %q, expected THREAT-001", threat.ID)
	}
	if threat.Category != CategoryTool {
		t.Errorf("Category = %q, expected tool", threat.Category)
	}
	if threat.Severity != SeverityHigh {
		t.Errorf("Severity = %q, expected high", threat.Severity)
	}
	if threat.Confidence != 0.90 {
		t.Errorf("Confidence = %f, expected 0.90", threat.Confidence)
	}
	if threat.Action != ActionBlock {
		t.Errorf("Action = %q, expected block", threat.Action)
	}
	if threat.Title != "SQL Injection" {
		t.Errorf("Title = %q, expected 'SQL Injection'", threat.Title)
	}
}

func TestParseDirectives(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Directive
	}{
		{
			name:     "empty content",
			content:  "",
			expected: nil,
		},
		{
			name:    "single BLOCK directive",
			content: `BLOCK: tool.call exec`,
			expected: []Directive{
				{Action: ActionBlock, Condition: "tool.call exec"},
			},
		},
		{
			name: "multiple directives",
			content: `BLOCK: tool.call exec
APPROVE: skill name contains untrusted
LOG: tool.call monitored`,
			expected: []Directive{
				{Action: ActionBlock, Condition: "tool.call exec"},
				{Action: ActionRequireApproval, Condition: "skill name contains untrusted"},
				{Action: ActionLog, Condition: "tool.call monitored"},
			},
		},
		{
			name: "case insensitive actions",
			content: `block: tool.call exec
Approve: skill name contains untrusted
log: tool.call monitored`,
			expected: []Directive{
				{Action: ActionBlock, Condition: "tool.call exec"},
				{Action: ActionRequireApproval, Condition: "skill name contains untrusted"},
				{Action: ActionLog, Condition: "tool.call monitored"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directives := ParseDirectives(tt.content)
			if len(directives) != len(tt.expected) {
				t.Errorf("ParseDirectives() returned %d directives, expected %d", len(directives), len(tt.expected))
				return
			}
			for i, d := range directives {
				if d.Action != tt.expected[i].Action {
					t.Errorf("Directive[%d].Action = %q, expected %q", i, d.Action, tt.expected[i].Action)
				}
				if d.Condition != tt.expected[i].Condition {
					t.Errorf("Directive[%d].Condition = %q, expected %q", i, d.Condition, tt.expected[i].Condition)
				}
			}
		})
	}
}

func TestIsValidThreat(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		threat  *ThreatEntry
		isValid bool
	}{
		{
			name: "valid threat",
			threat: &ThreatEntry{
				ID:         "THREAT-001",
				Category:   CategoryTool,
				Severity:   SeverityHigh,
				Confidence: 0.9,
				Action:     ActionBlock,
				Title:      "Test",
			},
			isValid: true,
		},
		{
			name: "missing ID",
			threat: &ThreatEntry{
				Category:   CategoryTool,
				Severity:   SeverityHigh,
				Confidence: 0.9,
				Action:     ActionBlock,
				Title:      "Test",
			},
			isValid: false,
		},
		{
			name: "invalid category",
			threat: &ThreatEntry{
				ID:         "THREAT-001",
				Category:   "invalid",
				Severity:   SeverityHigh,
				Confidence: 0.9,
				Action:     ActionBlock,
				Title:      "Test",
			},
			isValid: false,
		},
		{
			name: "invalid severity",
			threat: &ThreatEntry{
				ID:         "THREAT-001",
				Category:   CategoryTool,
				Severity:   "invalid",
				Confidence: 0.9,
				Action:     ActionBlock,
				Title:      "Test",
			},
			isValid: false,
		},
		{
			name: "invalid action",
			threat: &ThreatEntry{
				ID:         "THREAT-001",
				Category:   CategoryTool,
				Severity:   SeverityHigh,
				Confidence: 0.9,
				Action:     "invalid",
				Title:      "Test",
			},
			isValid: false,
		},
		{
			name: "confidence out of range",
			threat: &ThreatEntry{
				ID:         "THREAT-001",
				Category:   CategoryTool,
				Severity:   SeverityHigh,
				Confidence: 1.5,
				Action:     ActionBlock,
				Title:      "Test",
			},
			isValid: false,
		},
		{
			name: "missing title",
			threat: &ThreatEntry{
				ID:         "THREAT-001",
				Category:   CategoryTool,
				Severity:   SeverityHigh,
				Confidence: 0.9,
				Action:     ActionBlock,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.isValidThreat(tt.threat)
			if result != tt.isValid {
				t.Errorf("isValidThreat() = %v, expected %v", result, tt.isValid)
			}
		})
	}
}

func TestIsExpiredOrRevoked(t *testing.T) {
	parser := NewParser()

	// Test revoked threat
	revokedThreat := &ThreatEntry{Revoked: true}
	if !parser.isExpiredOrRevoked(revokedThreat) {
		t.Error("Revoked threat should be expired/revoked")
	}

	// Test expired threat
	pastTime := time.Now().Add(-24 * time.Hour)
	expiredThreat := &ThreatEntry{ExpiresAt: &pastTime}
	if !parser.isExpiredOrRevoked(expiredThreat) {
		t.Error("Expired threat should be expired/revoked")
	}

	// Test active threat
	futureTime := time.Now().Add(24 * time.Hour)
	activeThreat := &ThreatEntry{ExpiresAt: &futureTime}
	if parser.isExpiredOrRevoked(activeThreat) {
		t.Error("Active threat should not be expired/revoked")
	}

	// Test threat with no expiration
	noExpThreat := &ThreatEntry{}
	if parser.isExpiredOrRevoked(noExpThreat) {
		t.Error("Threat with no expiration should not be expired/revoked")
	}
}
