# go-countries

Country identifier resolver for Go. Convert between country codes, names, and emoji flags.

## Features

- ✅ Support for ISO 3166-1 alpha-2 codes (e.g., "VN")
- ✅ Support for ISO 3166-1 alpha-3 codes (e.g., "VNM")
- ✅ Support for CIOC codes (e.g., "GER")
- ✅ Support for country names and aliases (e.g., "Vietnam", "UK")
- ✅ Support for Great Britain subdivisions (England, Scotland, Wales)
- ✅ Fuzzy matching for typos and variations (Levenshtein distance ≤ 2)
- ✅ Flag emoji conversion

## Install
```
go get -u github.com/yudgnahk/go-countries
```

## Usage

### Get Country Information (Recommended)

The `GetCountryInfo` function is the main entry point. It accepts any input and returns all country information:

```go
package main

import (
	"fmt"

	countries "github.com/yudgnahk/go-countries"
)

func main() {
	// From country name
	info := countries.GetCountryInfo("Vietnam")
	fmt.Printf("Alpha2: %s, Alpha3: %s, CIOC: %s, Name: %s\n", info.Alpha2, info.Alpha3, info.CIOC, info.Name)
	// Output: Alpha2: VN, Alpha3: VNM, CIOC: VIE, Name: Vietnam

	// From alpha-3 code
	info = countries.GetCountryInfo("VNM")
	fmt.Printf("Alpha2: %s, Name: %s\n", info.Alpha2, info.Name)
	// Output: Alpha2: VN, Name: Vietnam

	// From flag emoji
	info = countries.GetCountryInfo("🇺🇸")
	fmt.Printf("Alpha2: %s, Name: %s\n", info.Alpha2, info.Name)
	// Output: Alpha2: US, Name: United States

	// Fuzzy matching handles typos
	info = countries.GetCountryInfo("GERM")
	fmt.Printf("Alpha2: %s, Name: %s\n", info.Alpha2, info.Name)
	// Output: Alpha2: DE, Name: Germany
}
```

### Get Flag Emoji

Convert country codes to flag emojis:

```go
flag := countries.GetFlag("VN")     // Returns "🇻🇳"
flag = countries.GetFlag("VNM")    // Returns "🇻🇳"
flag = countries.GetFlag("GER")    // Returns "🇩🇪"
flag = countries.GetFlag("GB-ENG") // Returns "🏴󠁧󠁢󠁥󠁮󠁧󠁿"
```

### Get Country Name

Get country names from codes or flags:

```go
name := countries.GetName("VN")     // Returns "Vietnam"
name = countries.GetName("VNM")    // Returns "Vietnam"
name = countries.GetName("GER")    // Returns "Germany"
name = countries.GetName("🇻🇳")     // Returns "Vietnam"
```

## API Reference

### `GetFlag(countryCode string) string`
Converts a country code to its emoji flag. Supports ISO 3166-1 alpha-2, alpha-3, CIOC codes, and special subdivisions.

**Parameters:**
- `countryCode` - 2-letter (VN), 3-letter (VNM, GER), or special codes (GB-ENG)

**Returns:** Flag emoji string, or empty string if not found

### `GetName(input string) string`
Converts a country code or flag emoji to its country name.

**Parameters:**
- `input` - Country code (alpha-2, alpha-3, CIOC) or flag emoji

**Returns:** Country name, or empty string if not found

### `GetCountryInfo(input string) CountryInfo`
Returns complete country information for any input. Supports codes, names, aliases, and flag emojis. Includes fuzzy matching for handling typos.

**Parameters:**
- `input` - Country code, name, alias, or flag emoji

**Returns:** `CountryInfo` struct with:
- `Alpha2` - ISO 3166-1 alpha-2 code (e.g., "VN")
- `Alpha3` - ISO 3166-1 alpha-3 code (e.g., "VNM")
- `CIOC` - CIOC code (e.g., "VIE")
- `Name` - Country name (e.g., "Vietnam")

## Supported Codes

### ISO 3166-1 Codes
- **Alpha-2**: 2-letter codes (e.g., VN, US, GB)
- **Alpha-3**: 3-letter codes (e.g., VNM, USA, GBR)
- **CIOC**: Olympic codes (e.g., GER for Germany, SUI for Switzerland)

### Special Subdivisions
- **England**: `GB-ENG` or `ENG`
- **Scotland**: `GB-SCT` or `SCT`
- **Wales**: `GB-WLS` or `WLS`

## Data Maps

The library provides access to several data maps:

- `CountryNames` - map[string]string: ISO alpha-2 codes to country names
- `CountryAliases` - map[string]string: Common aliases to ISO alpha-2 codes
- `Cca2CodeMap` - map[string]string: ISO alpha-2 code mappings
- `Cca3CodeMap` - map[string]string: ISO alpha-3 to alpha-2 code mappings
- `CiocCodeMap` - map[string]string: CIOC to alpha-2 code mappings
- `SpecialEmojiMap` - map[string]string: Special subdivision codes to emoji flags

### Special thanks to those repositories which helps me so much:
 - [go-emoji-flag](https://github.com/jayco/go-emoji-flag)
 - [restcountries](https://github.com/apilayer/restcountries)

... and wikipedia for the understanding about countries code (especially ISO_3166) (https://en.wikipedia.org/wiki/ISO_3166-2:GB)

I was stuck on the countries that belongs to Great Britain, which doesn't have the same format as origin emojis
([go-emoji-flag](https://github.com/jayco/go-emoji-flag) does not support these countries)
