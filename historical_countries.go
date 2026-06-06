package emojiflags

// HistoricalCountry represents a nation that no longer exists as an
// independent state but still has a valid (or formerly valid) ISO 3166
// code that can be resolved to a flag emoji sequence.
//
// Code references:
//   - ISO 3166-1 alpha-2: the 2-letter code assigned during the country's
//     existence. May be shared between countries that succeeded each other
//     (notably CS for both Czechoslovakia and Serbia and Montenegro).
//   - ISO 3166-1 alpha-3: the 3-letter code. Unique per historical country.
//   - IOC/CIOC: the Olympic Committee code. Unique per historical country.
//   - ISO 3166-3 alpha-4: the stable identifier assigned by ISO for
//     formerly used country names. Unique per historical country.
type HistoricalCountry struct {
	Alpha4 string // ISO 3166-3 alpha-4 (e.g. "SUHH", "CSHH", "CSXX")
	Alpha2 string // ISO 3166-1 alpha-2 (e.g. "SU", "CS")
	Alpha3 string // ISO 3166-1 alpha-3 (e.g. "SUN", "CSK", "SCG")
	Cioc   string // IOC/CIOC code (e.g. "URS", "TCH", "GDR")
	Name   string // Canonical English name
}

// historicalCountries lists nations that have participated in the FIFA
// World Cup finals but no longer exist as independent countries.
//
// Order matters for ambiguous alpha-2 codes: the FIRST entry with a given
// alpha-2 wins when the lookup cannot be disambiguated by alpha-3, CIOC,
// or alpha-4. Czechoslovakia is listed before Serbia and Montenegro so
// that bare "CS" resolves to Czechoslovakia (the earlier user of the code).
//
// The flag for each entry is generated as a pair of Regional Indicator
// symbols from the alpha-2 code. On stock iOS / Android / macOS / Windows
// these render as a generic "missing flag" placeholder because the codes
// are deprecated in CLDR; consumers that want proper historical flag
// glyphs should load a webfont such as BabelStone Flags.
var historicalCountries = []HistoricalCountry{
	{Alpha4: "SUHH", Alpha2: "SU", Alpha3: "SUN", Cioc: "URS", Name: "Soviet Union"},
	{Alpha4: "YUCS", Alpha2: "YU", Alpha3: "YUG", Cioc: "YUG", Name: "Yugoslavia"},
	{Alpha4: "CSHH", Alpha2: "CS", Alpha3: "CSK", Cioc: "TCH", Name: "Czechoslovakia"},
	{Alpha4: "CSXX", Alpha2: "CS", Alpha3: "SCG", Cioc: "SCG", Name: "Serbia and Montenegro"},
	{Alpha4: "DDDE", Alpha2: "DD", Alpha3: "DDR", Cioc: "GDR", Name: "East Germany"},
	{Alpha4: "ZRCD", Alpha2: "ZR", Alpha3: "ZAR", Cioc: "ZAI", Name: "Zaire"},
}

// HistoricalCountries returns a snapshot of all historical nations
// supported by the library. A defensive copy is returned so that
// external mutations cannot desynchronize the internal lookup indices
// (which are built once in init()).
//
// The slice order matches the underlying source-of-truth list and
// therefore also defines the disambiguation default for the ambiguous
// alpha-2 code "CS" (Czechoslovakia wins over Serbia and Montenegro).
func HistoricalCountries() []HistoricalCountry {
	out := make([]HistoricalCountry, len(historicalCountries))
	copy(out, historicalCountries)
	return out
}

// historicalByAlpha4 is the master index keyed by the unique ISO 3166-3
// alpha-4 code. All other lookup maps are derived from this.
var historicalByAlpha4 map[string]HistoricalCountry

// historicalByAlpha2 maps the 2-letter ISO 3166-1 code to a default
// alpha-4. For ambiguous codes (only CS in this set) the first writer
// wins, so bare "CS" resolves to Czechoslovakia.
var historicalByAlpha2 map[string]string

// historicalByAlpha3 maps the 3-letter ISO 3166-1 code to its alpha-4.
var historicalByAlpha3 map[string]string

// historicalByCioc maps the IOC/CIOC code to its alpha-4.
var historicalByCioc map[string]string

func init() {
	historicalByAlpha4 = make(map[string]HistoricalCountry, len(historicalCountries))
	historicalByAlpha2 = make(map[string]string, len(historicalCountries))
	historicalByAlpha3 = make(map[string]string, len(historicalCountries))
	historicalByCioc = make(map[string]string, len(historicalCountries))

	for _, hc := range historicalCountries {
		historicalByAlpha4[hc.Alpha4] = hc
		// First writer wins for alpha-2 so the ordering in
		// historicalCountries defines the disambiguation default.
		if _, exists := historicalByAlpha2[hc.Alpha2]; !exists {
			historicalByAlpha2[hc.Alpha2] = hc.Alpha4
		}
		historicalByAlpha3[hc.Alpha3] = hc.Alpha4
		if hc.Cioc != "" {
			historicalByCioc[hc.Cioc] = hc.Alpha4
		}
	}
}

// lookupHistorical resolves any supported code (alpha-2, alpha-3, CIOC, or
// alpha-4) to its HistoricalCountry record. The returned bool is false
// when the code does not match a known historical country.
func lookupHistorical(code string) (HistoricalCountry, bool) {
	switch len(code) {
	case 2:
		if a4, ok := historicalByAlpha2[code]; ok {
			hc, ok := historicalByAlpha4[a4]
			return hc, ok
		}
	case 3:
		if a4, ok := historicalByAlpha3[code]; ok {
			hc, ok := historicalByAlpha4[a4]
			return hc, ok
		}
		if a4, ok := historicalByCioc[code]; ok {
			hc, ok := historicalByAlpha4[a4]
			return hc, ok
		}
	case 4:
		if hc, ok := historicalByAlpha4[code]; ok {
			return hc, true
		}
	}
	return HistoricalCountry{}, false
}
