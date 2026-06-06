package emojiflags

import "strings"

// GetFlag converts a country code (ISO 3166-1 alpha-2, alpha-3, or CIOC) to its corresponding emoji flag.
// It supports 2-letter codes (e.g., "VN"), 3-letter codes (e.g., "VNM" or "GER"),
// and special subdivision codes (e.g., "GB-ENG" for England, "ENG" for England short code).
// Returns an empty string if the country code is not found.
//
// Example:
//
//	flag := countries.GetFlag("VN")     // Returns "🇻🇳"
//	flag = countries.GetFlag("VNM")    // Returns "🇻🇳"
//	flag = countries.GetFlag("GER")    // Returns "🇩🇪"
//	flag = countries.GetFlag("GB-ENG") // Returns "🏴󠁧󠁢󠁥󠁮󠁧󠁿"
func GetFlag(countryCode string) string {
	countryCode = normalizeUpper(countryCode)

	// Historical countries (USSR, Yugoslavia, Czechoslovakia, etc.)
	// are resolved by alpha-2, alpha-3, CIOC, or alpha-4 code.
	if hc, ok := lookupHistorical(countryCode); ok {
		return codeToFlag(hc.Alpha2)
	}

	switch len(countryCode) {
	case 2:
		if code, ok := Cca2CodeMap[countryCode]; ok {
			return codeToFlag(code)
		}
	case 3:
		if code, ok := Cca3CodeMap[countryCode]; ok {
			return codeToFlag(code)
		}
		if code, ok := CiocCodeMap[countryCode]; ok {
			return codeToFlag(code)
		}
		if flag, ok := SpecialEmojiMap[countryCode]; ok {
			return flag
		}
	case 6:
		if flag, ok := SpecialEmojiMap[countryCode]; ok {
			return flag
		}
	}

	return ""
}

// GetName converts a country code or flag emoji to its country name.
// Supports ISO 3166-1 alpha-2, alpha-3, CIOC codes, and flag emojis.
// Returns an empty string if not found.
//
// Example:
//
//	name := countries.GetName("VN")    // Returns "Vietnam"
//	name = countries.GetName("VNM")   // Returns "Vietnam"
//	name = countries.GetName("🇻🇳")    // Returns "Vietnam"
func GetName(input string) string {
	normalized := normalizeUpper(input)

	// Try historical countries first so that ambiguous codes such as
	// "CS" (Czechoslovakia vs. Serbia and Montenegro) can be resolved
	// by the more specific alpha-3, CIOC, or alpha-4 input.
	if hc, ok := lookupHistorical(normalized); ok {
		return hc.Name
	}

	// Try alpha-2 code
	if name, ok := CountryNames[normalized]; ok {
		return name
	}

	// Try alpha-3 or CIOC code
	if len(normalized) == 3 {
		if code, ok := Cca3CodeMap[normalized]; ok {
			if name, ok := CountryNames[code]; ok {
				return name
			}
		}
		if code, ok := CiocCodeMap[normalized]; ok {
			if name, ok := CountryNames[code]; ok {
				return name
			}
		}
	}

	// Try special subdivision codes
	if canonical, ok := SpecialCountryMap[normalized]; ok {
		if name, ok := CountryNames[canonical]; ok {
			return name
		}
	}

	// Try flag emoji (check trimmed input)
	trimmed := strings.TrimSpace(input)
	if isFlagEmoji(trimmed) || isSpecialFlag(trimmed) {
		alpha2 := lookupAlpha2ByFlag(trimmed)
		if alpha2 != "" {
			if name, ok := CountryNames[alpha2]; ok {
				return name
			}
			if canonical, ok := SpecialCountryMap[alpha2]; ok {
				if name, ok := CountryNames[canonical]; ok {
					return name
				}
			}
		}
	}

	return ""
}

// CountryInfo holds all country information returned by GetCountryInfo.
type CountryInfo struct {
	Alpha2 string // ISO 3166-1 alpha-2 code (e.g., "VN") or canonical subdivision code (e.g., "GB-ENG")
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
	normalized := normalizeUpper(input)
	if normalized == "" {
		return CountryInfo{}
	}

	// Try historical countries first so that ambiguous codes such as
	// "CS" (Czechoslovakia vs. Serbia and Montenegro) can be resolved
	// by the more specific alpha-3, CIOC, or alpha-4 input.
	if hc, ok := lookupHistorical(normalized); ok {
		return CountryInfo{
			Alpha2: hc.Alpha2,
			Alpha3: hc.Alpha3,
			CIOC:   hc.Cioc,
			Name:   hc.Name,
		}
	}

	// Try exact matches
	alpha2 := lookupAlpha2ByCode(normalized)
	if alpha2 == "" {
		alpha2 = lookupAlpha2ByName(normalized)
	}
	if alpha2 == "" {
		alpha2 = lookupAlpha2ByAlias(normalized)
	}

	// Try flag emoji (check trimmed input)
	if alpha2 == "" {
		trimmed := strings.TrimSpace(input)
		if isFlagEmoji(trimmed) || isSpecialFlag(trimmed) {
			alpha2 = lookupAlpha2ByFlag(trimmed)
		}
	}

	// Try special subdivision codes with canonicalization
	if alpha2 == "" {
		if canonical, ok := SpecialCountryMap[normalized]; ok {
			info := buildCountryInfo(canonical)
			if name, ok := CountryNames[canonical]; ok {
				info.Name = name
			}
			return info
		}
	}

	// Try fuzzy matching
	if alpha2 == "" {
		alpha2 = fuzzyMatchCode(normalized)
	}
	if alpha2 == "" {
		alpha2 = fuzzyMatchName(normalized)
	}

	if alpha2 == "" {
		return CountryInfo{}
	}

	return buildCountryInfo(alpha2)
}
