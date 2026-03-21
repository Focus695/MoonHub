package delegation

// Helper functions for tools

// toStringSlice converts a value to a string slice
func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	switch arr := v.(type) {
	case []string:
		return arr
	case []any:
		var result []string
		for _, item := range arr {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// getOptionalBool returns a boolean value or the default
func getOptionalBool(v any, defVal bool) bool {
	if v == nil {
		return defVal
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return defVal
}
