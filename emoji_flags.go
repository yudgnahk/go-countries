package emojiflags

import (
	"strings"
)

const (
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

// GetName converts a country code or flag emoji to its country name.
// Supports ISO 3166-1 alpha-2, alpha-3, CIOC codes, and flag emojis.
// Returns an empty string if not found.
//
// Example:
//
//	name := countries.GetName("VN")    // Returns "Vietnam"
//	name := countries.GetName("VNM")   // Returns "Vietnam"
//	name := countries.GetName("🇻🇳")    // Returns "Vietnam"
func GetName(input string) string {
	input = strings.ToUpper(input)

	// Try to get country name by alpha-2 code first
	if name, ok := CountryNames[input]; ok {
		return name
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

	// Try as flag emoji (two regional indicator symbols)
	if len(input) >= 8 {
		runes := []rune(input)
		if len(runes) >= 2 && runes[0] >= 0x1F1E6 && runes[0] <= 0x1F1FF && runes[1] >= 0x1F1E6 && runes[1] <= 0x1F1FF {
			char1 := rune('A') + (runes[0] - 0x1F1E6)
			char2 := rune('A') + (runes[1] - 0x1F1E6)
			code := string(char1) + string(char2)
			if name, ok := CountryNames[code]; ok {
				return name
			}
		}
	}

	return ""
}

// CountryInfo holds all country information returned by GetCountryInfo.
type CountryInfo struct {
	Alpha2 string // ISO 3166-1 alpha-2 code (e.g., "VN")
	Alpha3 string // ISO 3166-1 alpha-3 code (e.g., "VNM")
	CIOC   string // CIOC code (e.g., "VIE")
	Name   string // Country name (e.g., "Vietnam")
}

// GetCountryInfo returns all country information for any input.
// Supports alpha-2 codes, alpha-3 codes, CIOC codes, country names, aliases, and flag emojis.
// Includes fuzzy matching for handling typos (Levenshtein distance ≤ 2).
// Returns an empty CountryInfo if no match is found.
//
// Example:
//
//	info := countries.GetCountryInfo("Vietnam")
//	// Returns: CountryInfo{Alpha2: "VN", Alpha3: "VNM", CIOC: "VIE", Name: "Vietnam"}
//
//	info = countries.GetCountryInfo("VNM")
//	// Returns: CountryInfo{Alpha2: "VN", Alpha3: "VNM", CIOC: "VIE", Name: "Vietnam"}
//
//	info = countries.GetCountryInfo("🇺🇸")
//	// Returns: CountryInfo{Alpha2: "US", Alpha3: "USA", CIOC: "USA", Name: "United States"}
//
//	info = countries.GetCountryInfo("GERM")
//	// Returns: CountryInfo{Alpha2: "DE", Alpha3: "DEU", CIOC: "GER", Name: "Germany"} (fuzzy match)
func GetCountryInfo(input string) CountryInfo {
	input = strings.TrimSpace(input)
	if input == "" {
		return CountryInfo{}
	}

	normalized := strings.ToUpper(input)
	var alpha2 string

	// Exact match: alpha-2 code
	if _, ok := Cca2CodeMap[normalized]; ok {
		alpha2 = normalized
	}

	// Exact match: alpha-3 code
	if alpha2 == "" {
		if code, ok := Cca3CodeMap[normalized]; ok {
			alpha2 = code
		}
	}

	// Exact match: CIOC code
	if alpha2 == "" {
		if code, ok := CiocCodeMap[normalized]; ok {
			alpha2 = code
		}
	}

	// Exact match: special subdivision code (e.g., "GB-ENG", "ENG")
	if alpha2 == "" {
		if _, ok := SpecialEmojiMap[normalized]; ok {
			// For special codes, we return them as Alpha2 since they don't have standard alpha-2 codes
			info := CountryInfo{Alpha2: normalized}
			if name, ok := CountryNames[normalized]; ok {
				info.Name = name
			}
			return info
		}
	}

	// Exact match: flag emoji
	if alpha2 == "" && len(input) >= 8 {
		runes := []rune(input)
		if len(runes) >= 2 && runes[0] >= 0x1F1E6 && runes[0] <= 0x1F1FF && runes[1] >= 0x1F1E6 && runes[1] <= 0x1F1FF {
			char1 := rune('A') + (runes[0] - 0x1F1E6)
			char2 := rune('A') + (runes[1] - 0x1F1E6)
			code := string(char1) + string(char2)
			if _, ok := Cca2CodeMap[code]; ok {
				alpha2 = code
			}
		}
	}

	// Exact match: country name
	if alpha2 == "" {
		for code, name := range CountryNames {
			if strings.ToUpper(name) == normalized {
				alpha2 = code
				break
			}
		}
	}

	// Exact match: alias
	if alpha2 == "" {
		if code, ok := CountryAliases[normalized]; ok {
			alpha2 = code
		}
	}

	// Fuzzy match: codes (alpha-2, alpha-3, CIOC)
	if alpha2 == "" {
		const maxDistance = 2
		bestMatch := ""
		bestDistance := maxDistance + 1
		bestLength := 1000

		for code := range Cca2CodeMap {
			dist := levenshtein(normalized, code)
			if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
				bestDistance = dist
				bestMatch = code
				bestLength = len(code)
			}
		}

		for code := range Cca3CodeMap {
			dist := levenshtein(normalized, code)
			if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
				bestDistance = dist
				bestMatch = code
				bestLength = len(code)
			}
		}

		for code := range CiocCodeMap {
			dist := levenshtein(normalized, code)
			if dist < bestDistance || (dist == bestDistance && len(code) < bestLength) {
				bestDistance = dist
				bestMatch = code
				bestLength = len(code)
			}
		}

		if bestDistance <= maxDistance && bestMatch != "" {
			alpha2 = bestMatch
			// Normalize to alpha-2 if fuzzy matched on alpha-3 or CIOC
			if code, ok := Cca3CodeMap[bestMatch]; ok {
				alpha2 = code
			}
			if code, ok := CiocCodeMap[bestMatch]; ok {
				alpha2 = code
			}
		}
	}

	// Fuzzy match: country names and aliases
	if alpha2 == "" {
		const maxDistance = 2
		bestMatch := ""
		bestDistance := maxDistance + 1

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

		if bestDistance <= maxDistance && bestMatch != "" {
			alpha2 = bestMatch
		}
	}

	if alpha2 == "" {
		return CountryInfo{}
	}

	// Build result
	info := CountryInfo{Alpha2: alpha2}

	// Get name
	if name, ok := CountryNames[alpha2]; ok {
		info.Name = name
	}

	// Get alpha-3
	for code, cca2 := range Cca3CodeMap {
		if cca2 == alpha2 {
			info.Alpha3 = code
			break
		}
	}

	// Get CIOC
	for code, cca2 := range CiocCodeMap {
		if cca2 == alpha2 {
			info.CIOC = code
			break
		}
	}

	return info
}
