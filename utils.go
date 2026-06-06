package emojiflags

import "strings"

const (
	englandTagFlag  = "\U0001F3F4\U000E0067\U000E0062\U000E0065\U000E006E\U000E0067\U000E007F"
	scotlandTagFlag = "\U0001F3F4\U000E0067\U000E0062\U000E0073\U000E0063\U000E0074\U000E007F"
	walesTagFlag    = "\U0001F3F4\U000E0067\U000E0062\U000E0077\U000E006C\U000E0073\U000E007F"

	maxFuzzyDistance = 2
)

var SpecialEmojiMap = map[string]string{
	EnglandCode:      englandTagFlag,
	ScotlandCode:     scotlandTagFlag,
	WalesCode:        walesTagFlag,
	EnglandShortCode: englandTagFlag,
}

// specialFlagToCode provides reverse lookup for tag-sequence subdivision flags.
var specialFlagToCode map[string]string

func init() {
	specialFlagToCode = make(map[string]string, len(SpecialEmojiMap))
	for code, flag := range SpecialEmojiMap {
		specialFlagToCode[flag] = code
	}
}

// codeToFlag converts a 2-letter ISO code to its emoji flag representation.
func codeToFlag(code string) string {
	if len(code) != 2 {
		return ""
	}
	return string(0x1F1E6+rune(code[0])-'A') + string(0x1F1E6+rune(code[1])-'A')
}

// flagToCode converts an emoji flag to its 2-letter ISO code.
// Returns empty string if input is not a valid flag emoji.
func flagToCode(flag string) string {
	runes := []rune(flag)
	if len(runes) < 2 {
		return ""
	}

	// Check if both runes are regional indicator symbols (0x1F1E6-0x1F1FF)
	if runes[0] < 0x1F1E6 || runes[0] > 0x1F1FF || runes[1] < 0x1F1E6 || runes[1] > 0x1F1FF {
		return ""
	}

	char1 := rune('A') + (runes[0] - 0x1F1E6)
	char2 := rune('A') + (runes[1] - 0x1F1E6)
	return string(char1) + string(char2)
}

// isFlagEmoji checks if the input string is a flag emoji (two regional indicator symbols).
func isFlagEmoji(input string) bool {
	runes := []rune(input)
	if len(runes) < 2 {
		return false
	}
	return runes[0] >= 0x1F1E6 && runes[0] <= 0x1F1FF &&
		runes[1] >= 0x1F1E6 && runes[1] <= 0x1F1FF
}

// isSpecialFlag checks if the input string is a tag-sequence subdivision flag.
func isSpecialFlag(input string) bool {
	_, ok := specialFlagToCode[input]
	return ok
}

// normalizeUpper trims whitespace and converts to uppercase.
func normalizeUpper(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}

// lookupAlpha2ByCode tries to find alpha-2 code from various code maps.
func lookupAlpha2ByCode(normalized string) string {
	// Try alpha-2 code
	if _, ok := Cca2CodeMap[normalized]; ok {
		return normalized
	}

	// Try alpha-3 code
	if code, ok := Cca3CodeMap[normalized]; ok {
		return code
	}

	// Try CIOC code
	if code, ok := CiocCodeMap[normalized]; ok {
		return code
	}

	return ""
}

// lookupAlpha2ByName finds alpha-2 code from country name.
func lookupAlpha2ByName(normalized string) string {
	for code, name := range CountryNames {
		if strings.ToUpper(name) == normalized {
			return code
		}
	}
	return ""
}

// lookupAlpha2ByAlias finds alpha-2 code from alias map.
func lookupAlpha2ByAlias(normalized string) string {
	if code, ok := CountryAliases[normalized]; ok {
		return code
	}
	return ""
}

// lookupAlpha2ByFlag tries to find alpha-2 code from flag emoji (regional indicator or tag-sequence).
// Historical alpha-2 codes (SU, YU, CS, DD, ZR) are also recognised so that
// GetName / GetCountryInfo can round-trip a historical flag back to a name.
func lookupAlpha2ByFlag(input string) string {
	trimmed := strings.TrimSpace(input)

	// Check for tag-sequence subdivision flags first
	if isSpecialFlag(trimmed) {
		if code, ok := specialFlagToCode[trimmed]; ok {
			return code
		}
	}

	// Check for standard regional indicator flags
	if isFlagEmoji(trimmed) {
		code := flagToCode(trimmed)
		if _, ok := Cca2CodeMap[code]; ok {
			return code
		}
		if _, ok := historicalByAlpha2[code]; ok {
			return code
		}
	}

	return ""
}

// levenshtein calculates the Levenshtein distance between two strings using runes.
// This correctly handles multi-byte UTF-8 characters (e.g., "Curaçao", "Réunion").
func levenshtein(s1, s2 string) int {
	runes1 := []rune(s1)
	runes2 := []rune(s2)

	if len(runes1) == 0 {
		return len(runes2)
	}
	if len(runes2) == 0 {
		return len(runes1)
	}

	matrix := make([][]int, len(runes1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(runes2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(runes1); i++ {
		for j := 1; j <= len(runes2); j++ {
			cost := 1
			if runes1[i-1] == runes2[j-1] {
				cost = 0
			}

			deletion := matrix[i-1][j] + 1
			insertion := matrix[i][j-1] + 1
			substitution := matrix[i-1][j-1] + cost

			min := deletion
			if insertion < min {
				min = insertion
			}
			if substitution < min {
				min = substitution
			}

			matrix[i][j] = min
		}
	}

	return matrix[len(runes1)][len(runes2)]
}

// fuzzyMatchCode finds the closest matching code using Levenshtein distance.
func fuzzyMatchCode(normalized string) string {
	bestMatch := ""
	bestDistance := maxFuzzyDistance + 1
	bestLength := 1000

	// Check alpha-2 codes
	for code := range Cca2CodeMap {
		dist := levenshtein(normalized, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	// Check alpha-3 codes
	for code := range Cca3CodeMap {
		dist := levenshtein(normalized, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	// Check CIOC codes
	for code := range CiocCodeMap {
		dist := levenshtein(normalized, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	// Check special subdivision codes
	for code := range SpecialCountryMap {
		dist := levenshtein(normalized, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = SpecialCountryMap[code]
			bestLength = len(code)
		}
	}

	if bestDistance > maxFuzzyDistance || bestMatch == "" {
		return ""
	}

	// Normalize to alpha-2 if fuzzy matched on alpha-3 or CIOC
	if code, ok := Cca3CodeMap[bestMatch]; ok {
		return code
	}
	if code, ok := CiocCodeMap[bestMatch]; ok {
		return code
	}
	return bestMatch
}

// fuzzyMatchName finds the closest matching country name using Levenshtein distance.
func fuzzyMatchName(normalized string) string {
	bestMatch := ""
	bestDistance := maxFuzzyDistance + 1

	for code, name := range CountryNames {
		dist := levenshtein(normalized, strings.ToUpper(name))
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = code
		}
	}

	for alias, code := range CountryAliases {
		dist := levenshtein(normalized, strings.ToUpper(alias))
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = code
		}
	}

	if bestDistance > maxFuzzyDistance || bestMatch == "" {
		return ""
	}
	return bestMatch
}

// buildCountryInfo constructs a CountryInfo struct from an alpha-2 code.
func buildCountryInfo(alpha2 string) CountryInfo {
	info := CountryInfo{Alpha2: alpha2}

	if name, ok := CountryNames[alpha2]; ok {
		info.Name = name
	}

	for code, cca2 := range Cca3CodeMap {
		if cca2 == alpha2 {
			info.Alpha3 = code
			break
		}
	}

	for code, cca2 := range CiocCodeMap {
		if cca2 == alpha2 {
			info.CIOC = code
			break
		}
	}

	return info
}
