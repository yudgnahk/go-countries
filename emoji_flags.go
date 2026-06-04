package emojiflags

import (
	"strings"
)

const (
	ResolveKindCode      = "code"
	ResolveKindAlias     = "alias"
	ResolveKindName      = "name"
	ResolveKindFuzzyCode = "fuzzy-code"
	ResolveKindFuzzyName = "fuzzy-name"

	englandTagFlag  = "\U0001F3F4\U000E0067\U000E0062\U000E0065\U000E006E\U000E0067\U000E007F"
	scotlandTagFlag = "\U0001F3F4\U000E0067\U000E0062\U000E0073\U000E0063\U000E0074\U000E007F"
	walesTagFlag    = "\U0001F3F4\U000E0067\U000E0062\U000E0077\U000E006C\U000E0073\U000E007F"
)

var SpecialEmojiMap = map[string]string{
	EnglandCode:      englandTagFlag,
	ScotlandCode:     scotlandTagFlag,
	WalesCode:        walesTagFlag,
	EnglandShortCode: englandTagFlag,
}

// GetFlag converts a country code (ISO 3166-1 alpha-2, alpha-3, or CIOC) to its corresponding emoji flag.
// It supports 2-letter codes (e.g., "VN"), 3-letter codes (e.g., "VNM" or "GER"),
// and special subdivision codes (e.g., "GB-ENG" for England, "ENG" for England short code).
// Returns an empty string if the country code is not found.
func GetFlag(countryCode string) string {
	countryCode = strings.ToUpper(countryCode)
	switch len(countryCode) {
	case 2:
		if code, ok := Cca2CodeMap[countryCode]; ok {
			return string(0x1F1E6+rune(code[0])-'A') + string(0x1F1E6+rune(code[1])-'A')
		}
	case 3:
		if code, ok := Cca3CodeMap[countryCode]; ok {
			return string(0x1F1E6+rune(code[0])-'A') + string(0x1F1E6+rune(code[1])-'A')
		}

		if code, ok := CiocCodeMap[countryCode]; ok {
			return string(0x1F1E6+rune(code[0])-'A') + string(0x1F1E6+rune(code[1])-'A')
		}

		if flag, ok := SpecialEmojiMap[countryCode]; ok {
			return flag
		}
	case 6:
		if flag, ok := SpecialEmojiMap[countryCode]; ok {
			return flag
		}
	default:
		return ""
	}

	return ""
}

// GetFlagFuzzy attempts to find a flag using fuzzy matching on country codes.
// It searches for the closest match within all code maps (alpha-2, alpha-3, CIOC, and special codes).
// Returns the flag and the matched code if a close match is found (distance <= 2), otherwise returns empty strings.
// This is useful for handling typos or variations in country code input.
// When multiple codes have the same distance, it prefers shorter codes for more intuitive results.
//
// Example:
//
//	flag, code := emojiflags.GetFlagFuzzy("VIETNM")  // Returns Vietnam flag and "VNM"
//	flag, code := emojiflags.GetFlagFuzzy("USA")     // Also works when the input is already a valid code
func GetFlagFuzzy(input string) (string, string) {
	input = strings.ToUpper(input)

	// Try exact match first
	if flag := GetFlag(input); flag != "" {
		return flag, input
	}

	const maxDistance = 2
	bestMatch := ""
	bestDistance := maxDistance + 1
	bestLength := 1000 // Prefer shorter codes

	// Check alpha-2 codes
	for code := range Cca2CodeMap {
		dist := levenshtein(input, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	// Check alpha-3 codes
	for code := range Cca3CodeMap {
		dist := levenshtein(input, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	// Check CIOC codes
	for code := range CiocCodeMap {
		dist := levenshtein(input, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	// Check special codes
	for code := range SpecialEmojiMap {
		dist := levenshtein(input, code)
		if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
			bestDistance = dist
			bestMatch = code
			bestLength = len(code)
		}
	}

	if bestDistance <= maxDistance && bestMatch != "" {
		flag := GetFlag(bestMatch)
		return flag, bestMatch
	}

	return "", ""
}

// levenshtein calculates the Levenshtein distance between two strings.
// This measures the minimum number of single-character edits (insertions, deletions, or substitutions)
// required to change one string into the other.
func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
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

	return matrix[len(s1)][len(s2)]
}

// GetCode converts a flag emoji to its corresponding ISO 3166-1 alpha-2 country code.
// Returns an empty string if the flag is not recognized.
//
// Example:
//
//	code := emojiflags.GetCode("🇻🇳")  // Returns "VN"
//	code := emojiflags.GetCode("🏴󠁧󠁢󠁥󠁮󠁧󠁿")  // Returns "GB-ENG"
func GetCode(flag string) string {
	// Check special flags first.
	if flag == englandTagFlag {
		return EnglandCode
	}
	if flag == scotlandTagFlag {
		return ScotlandCode
	}
	if flag == walesTagFlag {
		return WalesCode
	}

	// The plain black flag emoji is shared across subdivisions,
	// so we use England as deterministic default.
	if flag == "🏴" {
		return EnglandCode
	}

	// Check if it's a standard flag emoji (two regional indicator symbols)
	if len(flag) < 8 {
		return ""
	}

	// Extract the two regional indicator symbols
	runes := []rune(flag)
	if len(runes) < 2 {
		return ""
	}

	// Regional indicator symbols are in range 0x1F1E6-0x1F1FF
	if runes[0] < 0x1F1E6 || runes[0] > 0x1F1FF || runes[1] < 0x1F1E6 || runes[1] > 0x1F1FF {
		return ""
	}

	// Convert regional indicators back to alpha-2 code
	char1 := rune('A') + (runes[0] - 0x1F1E6)
	char2 := rune('A') + (runes[1] - 0x1F1E6)
	code := string(char1) + string(char2)

	// Verify the code exists in our map
	if _, ok := Cca2CodeMap[code]; ok {
		return code
	}

	return ""
}

// ResolveFlag resolves mixed country identifiers to a flag.
// It supports exact code, exact alias/name, then fuzzy code/name matching.
// Returns the flag, matched code, and match kind.
func ResolveFlag(input string) (string, string, string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", ""
	}

	normalized := strings.ToUpper(input)

	if flag := GetFlag(normalized); flag != "" {
		return flag, normalized, ResolveKindCode
	}

	for code, countryName := range CountryNames {
		if strings.ToUpper(countryName) == normalized {
			return GetFlag(code), code, ResolveKindName
		}
	}

	if code, ok := CountryAliases[normalized]; ok {
		if flag := GetFlag(code); flag != "" {
			return flag, code, ResolveKindAlias
		}
	}

	if flag, code := GetFlagFuzzy(normalized); flag != "" {
		return flag, code, ResolveKindFuzzyCode
	}

	if flag, code := GetFlagByName(normalized); flag != "" {
		return flag, code, ResolveKindFuzzyName
	}

	return "", "", ""
}

// GetName converts a country code or flag emoji to its country name.
// Supports ISO 3166-1 alpha-2, alpha-3, CIOC codes, and flag emojis.
// Returns an empty string if not found.
//
// Example:
//
//	name := emojiflags.GetName("VN")    // Returns "Vietnam"
//	name := emojiflags.GetName("VNM")   // Returns "Vietnam"
//	name := emojiflags.GetName("🇻🇳")    // Returns "Vietnam"
func GetName(input string) string {
	input = strings.ToUpper(input)

	// Try to get country name by code first
	if name, ok := CountryNames[input]; ok {
		return name
	}

	// If input looks like a flag emoji, convert to code first
	if len(input) >= 8 {
		code := GetCode(input)
		if code != "" {
			if name, ok := CountryNames[code]; ok {
				return name
			}
		}
	}

	// Try alpha-3 codes
	if len(input) == 3 {
		if code, ok := Cca3CodeMap[input]; ok {
			if name, ok := CountryNames[code]; ok {
				return name
			}
		}
		// Try CIOC codes
		if code, ok := CiocCodeMap[input]; ok {
			if name, ok := CountryNames[code]; ok {
				return name
			}
		}
	}

	return ""
}

// GetFlagByName attempts to find a flag emoji by country name.
// Supports exact matches and fuzzy matching for country names and common aliases.
// Returns the flag emoji and matched code if found, empty strings otherwise.
//
// Example:
//
//	flag, code := emojiflags.GetFlagByName("Vietnam")     // Returns "🇻🇳", "VN"
//	flag, code := emojiflags.GetFlagByName("United States") // Returns "🇺🇸", "US"
func GetFlagByName(name string) (string, string) {
	name = strings.ToUpper(name)

	// Return empty for empty input
	if name == "" {
		return "", ""
	}

	// Try exact match first
	for code, countryName := range CountryNames {
		if strings.ToUpper(countryName) == name {
			return GetFlag(code), code
		}
	}

	if code, ok := CountryAliases[name]; ok {
		return GetFlag(code), code
	}

	// Try fuzzy matching with country names
	const maxDistance = 2
	bestMatch := ""
	bestDistance := maxDistance + 1

	for code, countryName := range CountryNames {
		dist := levenshtein(name, strings.ToUpper(countryName))
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = code
		}
	}

	// Also check common aliases in CountryAliases map
	for alias, code := range CountryAliases {
		dist := levenshtein(name, strings.ToUpper(alias))
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = code
		}
	}

	if bestDistance <= maxDistance && bestMatch != "" {
		return GetFlag(bestMatch), bestMatch
	}

	return "", ""
}

// GetCodes returns all country codes (alpha-2, alpha-3, CIOC) for a given country name.
// Returns empty strings for codes that are not found.
//
// Example:
//
//	alpha2, alpha3, cioc := emojiflags.GetCodes("Vietnam")  // Returns "VN", "VNM", "VIE"
//	alpha2, alpha3, cioc := emojiflags.GetCodes("Germany")  // Returns "DE", "DEU", "GER"
func GetCodes(name string) (string, string, string) {
	name = strings.ToUpper(name)
	if name == "" {
		return "", "", ""
	}

	// Find alpha-2 code from country name
	var alpha2 string
	for code, countryName := range CountryNames {
		if strings.ToUpper(countryName) == name {
			alpha2 = code
			break
		}
	}

	// Also check aliases
	if alpha2 == "" {
		if code, ok := CountryAliases[name]; ok {
			alpha2 = code
		}
	}

	if alpha2 == "" {
		return "", "", ""
	}

	// Find alpha-3 code from alpha-2
	var alpha3 string
	for code, cca2 := range Cca3CodeMap {
		if cca2 == alpha2 {
			alpha3 = code
			break
		}
	}

	// Find CIOC code from alpha-2
	var cioc string
	for code, cca2 := range CiocCodeMap {
		if cca2 == alpha2 {
			cioc = code
			break
		}
	}

	return alpha2, alpha3, cioc
}
