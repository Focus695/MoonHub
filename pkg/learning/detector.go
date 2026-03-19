// MoonHub - Ultra-lightweight personal AI agent
// Self-Improving Behavioral Pattern Detection System
// License: MIT

package learning

import (
	"regexp"
	"strings"
)

// PatternDetector detects behavioral patterns from conversations
type PatternDetector struct {
	config Config

	// Regex patterns for explicit feedback
	positivePatterns   []*regexp.Regexp
	negativePatterns   []*regexp.Regexp
	correctionPatterns []*regexp.Regexp

	// Semantic indicators
	semanticPositiveIndicators []SemanticIndicator
	semanticNegativeIndicators []SemanticIndicator

	// Workflow preference patterns
	workflowPatterns []*regexp.Regexp

	// Tool preference patterns
	toolPreferencePatterns []*regexp.Regexp
}

// SemanticIndicator represents a set of keywords with associated weight
type SemanticIndicator struct {
	Keywords []string
	Weight   float64
}

// DetectionContext contains the context for pattern detection
type DetectionContext struct {
	UserMessage         string
	AssistantMessage    string
	History             []Message
	ToolCalls           []ToolCall
	PreviousToolResults []ToolResult
}

// NewPatternDetector creates a new pattern detector
func NewPatternDetector(config Config) *PatternDetector {
	d := &PatternDetector{
		config: config,
	}

	d.initPatterns()
	d.initSemanticIndicators()

	return d
}

// initPatterns initializes regex patterns for detection
func (d *PatternDetector) initPatterns() {
	// Positive patterns - user acceptance signals
	d.positivePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(perfect|great|thanks|exactly|nice|awesome|good|yes|correct)`),
		regexp.MustCompile(`(?i)^(that'?s? (right|correct|it|what i (wanted|needed)))`),
		regexp.MustCompile(`(?i)^(love it|nailed it|spot on)`),
		regexp.MustCompile(`(?i)^(好的|很好|完美|谢谢|对|正确|正是我想要的)`), // Chinese
	}

	// Negative patterns - user rejection signals
	d.negativePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(no|wrong|not what|that'?s? not|you misunderstood|incorrect)`),
		regexp.MustCompile(`(?i)^(that'?s? (wrong|incorrect|not right))`),
		regexp.MustCompile(`(?i)^(try again|redo|not quite)`),
		regexp.MustCompile(`(?i)^(不对|错误|不是这个|再试一次|理解错了)`), // Chinese
	}

	// Correction patterns - user preference/correction signals
	d.correctionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(actually|i meant|i prefer|next time|please don'?t|instead)`),
		regexp.MustCompile(`(?i)^i (prefer|like|want|need) (.+)`),
		regexp.MustCompile(`(?i)^don'?t ([^,]+), (instead )?(.+)`),
		regexp.MustCompile(`(?i)^(remember|note) that i (.+)`),
		regexp.MustCompile(`(?i)^其实|我的意思是|我更喜欢|下次|不要|相反|记住我喜欢`), // Chinese
	}

	// Workflow preference patterns
	d.workflowPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)always (run|do|use|check) (.+) (before|after|when) (.+)`),
		regexp.MustCompile(`(?i)never (run|do|use) (.+) (when|if|unless) (.+)`),
		regexp.MustCompile(`(?i)(prefer|like) (it|to) (be|use|have) (.+)`),
		regexp.MustCompile(`(?i)每次|总是|从不|当.+时|之前|之后|喜欢用`), // Chinese
	}

	// Tool preference patterns
	d.toolPreferencePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)use (.+) instead of (.+)`),
		regexp.MustCompile(`(?i)don'?t use (.+)`),
		regexp.MustCompile(`(?i)(always|never) (run|call|use) (.+)`),
		regexp.MustCompile(`(?i)用.+代替|不要用|总是用|从不用`), // Chinese
	}
}

// initSemanticIndicators initializes semantic keyword indicators
func (d *PatternDetector) initSemanticIndicators() {
	d.semanticPositiveIndicators = []SemanticIndicator{
		{Keywords: []string{"helpful", "useful", "exactly", "perfect"}, Weight: 0.7},
		{Keywords: []string{"saved", "works", "solved", "fixed"}, Weight: 0.8},
		{Keywords: []string{"thanks", "appreciate", "great"}, Weight: 0.6},
		{Keywords: []string{"有帮助", "有用", "解决了", "谢谢", "太棒了"}, Weight: 0.7}, // Chinese
	}

	d.semanticNegativeIndicators = []SemanticIndicator{
		{Keywords: []string{"unhelpful", "wrong", "incorrect", "mistake"}, Weight: 0.8},
		{Keywords: []string{"not what", "didn't work", "failed", "error"}, Weight: 0.7},
		{Keywords: []string{"confusing", "unclear", "frustrating"}, Weight: 0.6},
		{Keywords: []string{"没帮助", "错了", "不正确", "失败", "困惑"}, Weight: 0.8}, // Chinese
	}
}

// DetectSignals detects all types of signals from the detection context
func (d *PatternDetector) DetectSignals(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal

	// 1. Regex-based detection
	signals = append(signals, d.detectRegexPatterns(ctx)...)

	// 2. Semantic analysis (if enabled)
	if d.config.EnableSemanticDetection {
		signals = append(signals, d.detectSemanticPatterns(ctx)...)
	}

	// 3. Workflow preferences
	signals = append(signals, d.detectWorkflowPreferences(ctx)...)

	// 4. Tool preferences
	signals = append(signals, d.detectToolPreferences(ctx)...)

	// 5. Implicit signals from tool usage (if enabled)
	if d.config.EnableImplicitSignals {
		signals = append(signals, d.detectImplicitToolSignals(ctx)...)
	}

	// 6. Behavioral patterns from conversation flow
	signals = append(signals, d.detectConversationFlowPatterns(ctx)...)

	return deduplicateSignals(signals)
}

// detectRegexPatterns detects patterns using regex
func (d *PatternDetector) detectRegexPatterns(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal
	lower := strings.ToLower(strings.TrimSpace(ctx.UserMessage))

	// Check positive patterns
	for _, pattern := range d.positivePatterns {
		if pattern.MatchString(lower) {
			signals = append(signals, DetectedSignal{
				Type:       CategorySuccessIndicator,
				Confidence: ConfidencePositive,
				Source:     SourceExplicitFeedback,
				Category:   CategorySuccessIndicator,
				Context:    truncateContext(ctx.AssistantMessage, 100),
				Timestamp:  NowMs(),
			})
		}
	}

	// Check negative patterns
	for _, pattern := range d.negativePatterns {
		if pattern.MatchString(lower) {
			signals = append(signals, DetectedSignal{
				Type:       CategoryFailureIndicator,
				Confidence: ConfidenceNegative,
				Source:     SourceExplicitFeedback,
				Category:   CategoryFailureIndicator,
				Context:    truncateContext(ctx.AssistantMessage, 100),
				Timestamp:  NowMs(),
			})
		}
	}

	// Check correction patterns
	for _, pattern := range d.correctionPatterns {
		if pattern.MatchString(lower) {
			signals = append(signals, DetectedSignal{
				Type:             CategoryCorrectionExplicit,
				Confidence:       ConfidenceCorrection,
				Source:           SourceExplicitFeedback,
				Category:         CategoryCorrectionExplicit,
				Context:          ctx.UserMessage,
				ExtractedPattern: extractPreference(ctx.UserMessage),
				Timestamp:        NowMs(),
			})
		}
	}

	return signals
}

// detectSemanticPatterns detects patterns using semantic keyword analysis
func (d *PatternDetector) detectSemanticPatterns(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal
	lower := strings.ToLower(ctx.UserMessage)

	// Check positive semantic indicators
	for _, indicator := range d.semanticPositiveIndicators {
		matchCount := 0
		for _, keyword := range indicator.Keywords {
			if strings.Contains(lower, keyword) {
				matchCount++
			}
		}
		if matchCount >= 2 {
			confidence := indicator.Weight * (float64(matchCount) / float64(len(indicator.Keywords)))
			signals = append(signals, DetectedSignal{
				Type:       CategorySuccessIndicator,
				Confidence: confidence,
				Source:     SourceExplicitFeedback,
				Category:   CategorySuccessIndicator,
				Context:    "Semantic positive: " + strings.Join(indicator.Keywords[:MinInt(3, len(indicator.Keywords))], ", "),
				Timestamp:  NowMs(),
			})
		}
	}

	// Check negative semantic indicators
	for _, indicator := range d.semanticNegativeIndicators {
		matchCount := 0
		for _, keyword := range indicator.Keywords {
			if strings.Contains(lower, keyword) {
				matchCount++
			}
		}
		if matchCount >= 2 {
			confidence := indicator.Weight * (float64(matchCount) / float64(len(indicator.Keywords)))
			signals = append(signals, DetectedSignal{
				Type:       CategoryFailureIndicator,
				Confidence: confidence,
				Source:     SourceExplicitFeedback,
				Category:   CategoryFailureIndicator,
				Context:    "Semantic negative: " + strings.Join(indicator.Keywords[:MinInt(3, len(indicator.Keywords))], ", "),
				Timestamp:  NowMs(),
			})
		}
	}

	return signals
}

// detectWorkflowPreferences detects workflow preference patterns
func (d *PatternDetector) detectWorkflowPreferences(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal
	lower := strings.ToLower(ctx.UserMessage)

	for _, pattern := range d.workflowPatterns {
		if pattern.MatchString(lower) {
			signals = append(signals, DetectedSignal{
				Type:             CategoryPreferenceWorkflow,
				Confidence:       0.85,
				Source:           SourceExplicitFeedback,
				Category:         CategoryPreferenceWorkflow,
				Context:          ctx.UserMessage,
				ExtractedPattern: ctx.UserMessage,
				Timestamp:        NowMs(),
			})
		}
	}

	return signals
}

// detectToolPreferences detects tool preference patterns
func (d *PatternDetector) detectToolPreferences(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal
	lower := strings.ToLower(ctx.UserMessage)

	for _, pattern := range d.toolPreferencePatterns {
		if pattern.MatchString(lower) {
			signals = append(signals, DetectedSignal{
				Type:             CategoryPreferenceTool,
				Confidence:       0.85,
				Source:           SourceExplicitFeedback,
				Category:         CategoryPreferenceTool,
				Context:          ctx.UserMessage,
				ExtractedPattern: ctx.UserMessage,
				Timestamp:        NowMs(),
			})
		}
	}

	return signals
}

// detectImplicitToolSignals detects implicit signals from tool usage
func (d *PatternDetector) detectImplicitToolSignals(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal

	if len(ctx.PreviousToolResults) == 0 {
		return signals
	}

	lower := strings.ToLower(strings.TrimSpace(ctx.UserMessage))

	// Acceptance indicators
	acceptanceIndicators := []string{"ok", "good", "thanks", "perfect", "great", "yes", "continue", "好的", "谢谢", "继续"}
	isAcceptance := false
	for _, indicator := range acceptanceIndicators {
		if lower == indicator || strings.HasPrefix(lower, indicator+" ") {
			isAcceptance = true
			break
		}
	}

	if isAcceptance {
		for _, result := range ctx.PreviousToolResults {
			if result.Success {
				signals = append(signals, DetectedSignal{
					Type:             CategorySuccessIndicator,
					Confidence:       ConfidenceImplicit,
					Source:           SourceToolUsageSuccess,
					Category:         CategorySuccessIndicator,
					Context:          "Tool " + result.Tool + " accepted",
					ExtractedPattern: "Tool " + result.Tool + " is appropriate in this context",
					Timestamp:        NowMs(),
				})
			}
		}
	}

	// Frustration indicators
	frustrationIndicators := []string{"try again", "that's not", "no", "wrong", "different", "再试", "不对", "错了"}
	isFrustration := false
	for _, indicator := range frustrationIndicators {
		if strings.Contains(lower, indicator) {
			isFrustration = true
			break
		}
	}

	if isFrustration {
		for _, result := range ctx.PreviousToolResults {
			signals = append(signals, DetectedSignal{
				Type:             CategoryFailureIndicator,
				Confidence:       0.6,
				Source:           SourceToolUsageFailure,
				Category:         CategoryFailureIndicator,
				Context:          "Tool " + result.Tool + " result was unsatisfactory",
				ExtractedPattern: "Tool " + result.Tool + " may not be optimal in this context",
				Timestamp:        NowMs(),
			})
		}
	}

	return signals
}

// detectConversationFlowPatterns detects patterns from conversation flow
func (d *PatternDetector) detectConversationFlowPatterns(ctx DetectionContext) []DetectedSignal {
	var signals []DetectedSignal

	// Analyze message length patterns
	if len(ctx.History) >= 3 {
		userMsgLen := len(ctx.UserMessage)
		shortResponses := 0
		for i := len(ctx.History) - 1; i >= 0 && i >= len(ctx.History)-3; i-- {
			if ctx.History[i].Role == "user" && len(ctx.History[i].Content) < 50 {
				shortResponses++
			}
		}

		// If user consistently gives short responses after assistant messages
		if shortResponses >= 2 && userMsgLen < 50 {
			signals = append(signals, DetectedSignal{
				Type:             CategoryPreferenceCommunication,
				Confidence:       0.5,
				Source:           SourceConversationFlow,
				Category:         CategoryPreferenceCommunication,
				Context:          "User prefers brief responses",
				ExtractedPattern: "Keep responses concise",
				Timestamp:        NowMs(),
			})
		}
	}

	return signals
}

// Helper functions

// extractPreference extracts a preference statement from a correction message
func extractPreference(message string) string {
	// Simple extraction - return the message for MVP
	// In a full implementation, this would use NLP to extract the core preference
	return message
}

// truncateContext truncates a context string to a maximum length
func truncateContext(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// deduplicateSignals removes duplicate signals
func deduplicateSignals(signals []DetectedSignal) []DetectedSignal {
	if signals == nil {
		return []DetectedSignal{}
	}
	seen := make(map[string]bool)
	var result []DetectedSignal

	for _, signal := range signals {
		key := string(signal.Type) + signal.Context
		if !seen[key] {
			seen[key] = true
			result = append(result, signal)
		}
	}

	return result
}
