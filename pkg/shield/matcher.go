package shield

import (
	"fmt"
	"regexp"
	"strings"
)

// Matcher evaluates events against threat conditions.
type Matcher struct{}

// NewMatcher creates a new condition matcher.
func NewMatcher() *Matcher {
	return &Matcher{}
}

// MatchEvent evaluates an event against all threats and returns matching results.
func (m *Matcher) MatchEvent(event ShieldEvent, threats []ThreatEntry) []MatchResult {
	var results []MatchResult

	for _, threat := range threats {
		// Skip if scope is not compatible with category
		if !m.isScopeCompatible(event.Scope, threat.Category) {
			continue
		}

		// Parse directives from recommendation_agent
		directives := ParseDirectives(threat.RecommendationAgent)
		if len(directives) == 0 {
			continue
		}

		// Evaluate each directive
		for _, directive := range directives {
			if matchInfo := m.evaluateCondition(directive.Condition, event); matchInfo != nil {
				results = append(results, MatchResult{
					Threat:     threat,
					Directive:  directive,
					MatchedOn:  matchInfo.matchedOn,
					MatchValue: matchInfo.matchValue,
				})
			}
		}
	}

	return results
}

// matchInfo holds information about a successful match.
type matchInfo struct {
	matchedOn  string
	matchValue string
}

// evaluateCondition evaluates a single condition against an event.
func (m *Matcher) evaluateCondition(condition string, event ShieldEvent) *matchInfo {
	condition = strings.TrimSpace(condition)
	if condition == "" {
		return nil
	}

	// Lowercase for case-insensitive matching
	condLower := strings.ToLower(condition)

	// Pattern: tool.call with arguments containing (...)
	// Must check this BEFORE "tool.call " pattern to avoid false matches
	if strings.HasPrefix(condLower, "tool.call with arguments containing ") {
		patternsStr := strings.TrimSpace(condition[36:])
		return m.matchToolArgsContaining(patternsStr, event)
	}

	// Pattern: tool.call <tool_name>
	if strings.HasPrefix(condLower, "tool.call ") {
		toolPattern := strings.TrimSpace(condition[10:])
		return m.matchToolCall(toolPattern, event)
	}

	// Pattern: skill name equals <value>
	if strings.HasPrefix(condLower, "skill name equals ") {
		value := strings.TrimSpace(condition[18:])
		return m.matchSkillNameEquals(value, event)
	}

	// Pattern: skill name contains <value>
	if strings.HasPrefix(condLower, "skill name contains ") {
		value := strings.TrimSpace(condition[20:])
		return m.matchSkillNameContains(value, event)
	}

	// Pattern: outbound request to <domain>
	if strings.HasPrefix(condLower, "outbound request to ") {
		domainPattern := strings.TrimSpace(condition[20:])
		return m.matchOutboundRequest(domainPattern, event)
	}

	// Pattern: file path equals <path>
	if strings.HasPrefix(condLower, "file path equals ") {
		path := strings.TrimSpace(condition[17:])
		return m.matchFilePathEquals(path, event)
	}

	// Pattern: file path contains <path>
	if strings.HasPrefix(condLower, "file path contains ") {
		path := strings.TrimSpace(condition[20:])
		return m.matchFilePathContains(path, event)
	}

	// Pattern: secrets read path equals <path>
	if strings.HasPrefix(condLower, "secrets read path equals ") {
		path := strings.TrimSpace(condition[25:])
		return m.matchSecretsPathEquals(path, event)
	}

	// Pattern: incoming message contains <text>
	if strings.HasPrefix(condLower, "incoming message contains ") {
		text := strings.TrimSpace(condition[26:])
		return m.matchIncomingMessageContains(text, event)
	}

	// Pattern: message contains <text>
	if strings.HasPrefix(condLower, "message contains ") {
		text := strings.TrimSpace(condition[17:])
		return m.matchIncomingMessageContains(text, event)
	}

	return nil
}

// matchToolCall matches tool.call <tool_name> patterns.
func (m *Matcher) matchToolCall(toolPattern string, event ShieldEvent) *matchInfo {
	if event.ToolName == "" {
		return nil
	}

	// Exact match or substring match
	toolPattern = strings.ToLower(toolPattern)
	toolNameLower := strings.ToLower(event.ToolName)

	if toolNameLower == toolPattern || strings.Contains(toolNameLower, toolPattern) {
		return &matchInfo{
			matchedOn:  "tool_name",
			matchValue: event.ToolName,
		}
	}

	return nil
}

// matchToolArgsContaining matches tool.call with arguments containing (...) patterns.
func (m *Matcher) matchToolArgsContaining(patternsStr string, event ShieldEvent) *matchInfo {
	if len(event.ToolArgs) == 0 {
		return nil
	}

	// Extract patterns from parentheses: (pattern1, pattern2, ...)
	patterns := extractPatternsFromParens(patternsStr)
	if len(patterns) == 0 {
		return nil
	}

	// Convert tool args to string for searching
	argsStr := argsToString(event.ToolArgs)
	argsStrLower := strings.ToLower(argsStr)

	// Check if any pattern is found in args
	for _, pattern := range patterns {
		patternLower := strings.ToLower(pattern)
		if strings.Contains(argsStrLower, patternLower) {
			return &matchInfo{
				matchedOn:  "tool_args",
				matchValue: pattern,
			}
		}
	}

	return nil
}

// matchSkillNameEquals matches skill name equals <value> patterns.
func (m *Matcher) matchSkillNameEquals(value string, event ShieldEvent) *matchInfo {
	if event.SkillName == "" {
		return nil
	}

	if strings.EqualFold(event.SkillName, value) {
		return &matchInfo{
			matchedOn:  "skill_name",
			matchValue: event.SkillName,
		}
	}

	return nil
}

// matchSkillNameContains matches skill name contains <value> patterns.
func (m *Matcher) matchSkillNameContains(value string, event ShieldEvent) *matchInfo {
	if event.SkillName == "" {
		return nil
	}

	if strings.Contains(strings.ToLower(event.SkillName), strings.ToLower(value)) {
		return &matchInfo{
			matchedOn:  "skill_name",
			matchValue: value,
		}
	}

	return nil
}

// matchOutboundRequest matches outbound request to <domain> patterns.
func (m *Matcher) matchOutboundRequest(domainPattern string, event ShieldEvent) *matchInfo {
	if event.Domain == "" {
		return nil
	}

	// Support OR operator: "domain1 or domain2"
	patterns := splitByOr(domainPattern)
	domainLower := strings.ToLower(event.Domain)

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if strings.Contains(domainLower, strings.ToLower(pattern)) {
			return &matchInfo{
				matchedOn:  "domain",
				matchValue: pattern,
			}
		}
	}

	return nil
}

// matchFilePathEquals matches file path equals <path> patterns.
func (m *Matcher) matchFilePathEquals(path string, event ShieldEvent) *matchInfo {
	if event.FilePath == "" {
		return nil
	}

	// Also check tool args for file paths
	filePath := event.FilePath
	if filePath == "" {
		filePath = getFilePathFromArgs(event.ToolArgs)
	}

	if filePath != "" && strings.EqualFold(filePath, path) {
		return &matchInfo{
			matchedOn:  "file_path",
			matchValue: filePath,
		}
	}

	return nil
}

// matchFilePathContains matches file path contains <path> patterns.
func (m *Matcher) matchFilePathContains(path string, event ShieldEvent) *matchInfo {
	// Check event.FilePath first
	if event.FilePath != "" {
		if strings.Contains(strings.ToLower(event.FilePath), strings.ToLower(path)) {
			return &matchInfo{
				matchedOn:  "file_path",
				matchValue: path,
			}
		}
	}

	// Also check tool args for file paths
	filePath := getFilePathFromArgs(event.ToolArgs)
	if filePath != "" {
		if strings.Contains(strings.ToLower(filePath), strings.ToLower(path)) {
			return &matchInfo{
				matchedOn:  "file_path",
				matchValue: path,
			}
		}
	}

	return nil
}

// matchSecretsPathEquals matches secrets read path equals <path> patterns.
func (m *Matcher) matchSecretsPathEquals(path string, event ShieldEvent) *matchInfo {
	if event.SecretPath == "" {
		return nil
	}

	// Support wildcards
	if strings.Contains(path, "*") {
		// Convert wildcard to regex
		regexPattern := "^" + regexp.QuoteMeta(path) + "$"
		regexPattern = strings.ReplaceAll(regexPattern, "\\*", ".*")
		matched, _ := regexp.MatchString(regexPattern, event.SecretPath)
		if matched {
			return &matchInfo{
				matchedOn:  "secret_path",
				matchValue: event.SecretPath,
			}
		}
	} else if strings.EqualFold(event.SecretPath, path) {
		return &matchInfo{
			matchedOn:  "secret_path",
			matchValue: event.SecretPath,
		}
	}

	return nil
}

// matchIncomingMessageContains matches incoming message contains <text> patterns.
func (m *Matcher) matchIncomingMessageContains(text string, event ShieldEvent) *matchInfo {
	if event.InputText == "" {
		return nil
	}

	if strings.Contains(strings.ToLower(event.InputText), strings.ToLower(text)) {
		return &matchInfo{
			matchedOn:  "input_text",
			matchValue: text,
		}
	}

	return nil
}

// isScopeCompatible checks if the event scope is compatible with the threat category.
func (m *Matcher) isScopeCompatible(scope ShieldScope, category ThreatCategory) bool {
	compatible := map[ShieldScope]map[ThreatCategory]bool{
		ScopeToolCall: {
			CategoryTool:   true,
			CategoryPrompt: true,
			CategoryMemory: true,
			CategoryFile:   true,
		},
		ScopeSkillInstall: {
			CategorySkill: true,
		},
		ScopeSkillExecute: {
			CategorySkill: true,
			CategoryTool:  true,
		},
		ScopeNetworkEgress: {
			CategoryNetwork: true,
		},
		ScopeSecretsRead: {
			CategoryTool:  true,
			CategoryFile:  true,
			CategoryOther: true,
		},
		ScopePrompt: {
			CategoryPrompt:  true,
			CategoryMessage: true,
		},
	}

	if categories, ok := compatible[scope]; ok {
		return categories[category]
	}
	return false
}

// extractPatternsFromParens extracts patterns from parentheses.
// Example: "(DROP, DELETE, UNION)" -> ["DROP", "DELETE", "UNION"]
func extractPatternsFromParens(s string) []string {
	// Find content in parentheses
	re := regexp.MustCompile(`\(([^)]+)\)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		// No parentheses, treat entire string as single pattern
		if s != "" {
			return []string{s}
		}
		return nil
	}

	// Split by comma and clean up
	parts := strings.Split(match[1], ",")
	patterns := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			patterns = append(patterns, part)
		}
	}
	return patterns
}

// splitByOr splits a string by " or " (case-insensitive).
func splitByOr(s string) []string {
	// Case-insensitive split by " or "
	re := regexp.MustCompile(`(?i)\s+or\s+`)
	parts := re.Split(s, -1)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// argsToString converts tool args map to a string for pattern matching.
func argsToString(args map[string]any) string {
	var sb strings.Builder
	for key, value := range args {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(fmt.Sprintf("%v", value))
		sb.WriteString(" ")
	}
	return sb.String()
}

// getFilePathFromArgs extracts file path from tool arguments.
func getFilePathFromArgs(args map[string]any) string {
	// Common argument names for file paths
	pathKeys := []string{"path", "file_path", "filepath", "file", "filename", "target", "destination"}

	for _, key := range pathKeys {
		if val, ok := args[key]; ok {
			if str, ok := val.(string); ok && str != "" {
				return str
			}
		}
	}

	return ""
}
