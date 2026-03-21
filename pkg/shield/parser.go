package shield

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser parses threat definitions from SHIELD.md content.
type Parser struct{}

// NewParser creates a new threat parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses SHIELD.md content and returns active threat entries.
func (p *Parser) Parse(content string) []ThreatEntry {
	if content == "" {
		return nil
	}

	// Find all YAML code blocks
	blocks := extractYAMLBlocks(content)
	if len(blocks) == 0 {
		return nil
	}

	var threats []ThreatEntry
	for _, block := range blocks {
		// Only process blocks that look like threat definitions
		if !strings.Contains(block, "id:") || !strings.Contains(block, "THREAT-") {
			continue
		}

		threat := p.parseThreatBlock(block)
		if threat != nil && p.isValidThreat(threat) && !p.isExpiredOrRevoked(threat) {
			threats = append(threats, *threat)
		}
	}

	return threats
}

// extractYAMLBlocks extracts YAML code blocks from markdown content.
func extractYAMLBlocks(content string) []string {
	// Match fenced code blocks with optional language specifier
	re := regexp.MustCompile("(?s)```(?:yaml|yml)?\n(.*?)```")
	matches := re.FindAllStringSubmatch(content, -1)

	blocks := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			blocks = append(blocks, strings.TrimSpace(match[1]))
		}
	}
	return blocks
}

// parseThreatBlock parses a single YAML threat block.
func (p *Parser) parseThreatBlock(block string) *ThreatEntry {
	threat := &ThreatEntry{}

	// Extract simple fields
	threat.ID = extractField(block, "id")
	threat.Fingerprint = extractField(block, "fingerprint")
	threat.Category = ThreatCategory(extractField(block, "category"))
	threat.Severity = ThreatSeverity(extractField(block, "severity"))
	threat.Action = ShieldAction(extractField(block, "action"))
	threat.Title = extractField(block, "title")
	threat.Description = extractField(block, "description")

	// Extract confidence
	if confStr := extractField(block, "confidence"); confStr != "" {
		if conf, err := strconv.ParseFloat(confStr, 64); err == nil {
			threat.Confidence = conf
		}
	}

	// Extract expires_at
	if expiresStr := extractField(block, "expires_at"); expiresStr != "" {
		if t, err := time.Parse(time.RFC3339, expiresStr); err == nil {
			threat.ExpiresAt = &t
		}
	}

	// Extract revoked
	if revokedStr := extractField(block, "revoked"); revokedStr != "" {
		threat.Revoked = strings.ToLower(revokedStr) == "true"
	}

	// Extract revoked_at
	if revokedAtStr := extractField(block, "revoked_at"); revokedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, revokedAtStr); err == nil {
			threat.RevokedAt = &t
		}
	}

	// Extract multi-line recommendation_agent
	threat.RecommendationAgent = extractMultilineField(block, "recommendation_agent")

	return threat
}

// extractField extracts a single-line YAML field value.
func extractField(block, field string) string {
	// Match: field: value or field: "value" or field: 'value'
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^%s:\s*["']?(.+?)["']?\s*$`, regexp.QuoteMeta(field)))
	match := re.FindStringSubmatch(block)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// extractMultilineField extracts a multi-line YAML field value (after | or >).
func extractMultilineField(block, field string) string {
	// Match: field: | followed by indented content
	re := regexp.MustCompile(fmt.Sprintf(`(?ms)^%s:\s*[|>]\s*\n((?:  .+\n?)+)`, regexp.QuoteMeta(field)))
	match := re.FindStringSubmatch(block)
	if len(match) > 1 {
		// Remove the leading spaces from each line
		lines := strings.Split(match[1], "\n")
		var result strings.Builder
		for _, line := range lines {
			// Remove exactly 2 spaces of indentation
			if strings.HasPrefix(line, "  ") {
				result.WriteString(line[2:])
			} else if strings.TrimSpace(line) != "" {
				result.WriteString(line)
			}
			if strings.TrimSpace(line) != "" {
				result.WriteString("\n")
			}
		}
		return strings.TrimSpace(result.String())
	}
	return ""
}

// isValidThreat validates that a threat entry has all required fields with valid values.
func (p *Parser) isValidThreat(threat *ThreatEntry) bool {
	if threat.ID == "" {
		return false
	}
	if !ValidCategories()[threat.Category] {
		return false
	}
	if !ValidSeverities()[threat.Severity] {
		return false
	}
	if !ValidActions()[threat.Action] {
		return false
	}
	if threat.Confidence < 0 || threat.Confidence > 1 {
		return false
	}
	if threat.Title == "" {
		return false
	}
	return true
}

// isExpiredOrRevoked checks if a threat is expired or revoked.
func (p *Parser) isExpiredOrRevoked(threat *ThreatEntry) bool {
	// Check if revoked
	if threat.Revoked {
		return true
	}

	// Check if expired
	if threat.ExpiresAt != nil && time.Now().After(*threat.ExpiresAt) {
		return true
	}

	return false
}

// ParseDirectives parses the recommendation_agent field into individual directives.
func ParseDirectives(recommendationAgent string) []Directive {
	if recommendationAgent == "" {
		return nil
	}

	var directives []Directive
	lines := strings.Split(recommendationAgent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse action prefixes: BLOCK:, APPROVE:, LOG:
		var action ShieldAction
		var condition string

		upperLine := strings.ToUpper(line)
		if strings.HasPrefix(upperLine, "BLOCK:") {
			action = ActionBlock
			condition = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(upperLine, "APPROVE:") {
			action = ActionRequireApproval
			condition = strings.TrimSpace(line[8:])
		} else if strings.HasPrefix(upperLine, "LOG:") {
			action = ActionLog
			condition = strings.TrimSpace(line[4:])
		} else {
			// Unknown format, skip
			continue
		}

		if condition != "" {
			directives = append(directives, Directive{
				Action:    action,
				Condition: condition,
			})
		}
	}

	return directives
}
