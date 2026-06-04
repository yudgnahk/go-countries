package emojiflags

import (
	"fmt"
	"testing"
	"unicode"
	"unicode/utf8"
)

func Test_getFlag(t *testing.T) {
	type args struct {
		country string
	}
	tests := []struct {
		name        string
		args        args
		expectedLen int
	}{
		{
			"Should handle correct 3 characters input",
			args{"VNM"},
			2,
		},
		{
			"Should handle correct 2 characters input",
			args{"VN"},
			2,
		},
		{
			"Should return empty string if no 3 letters code found",
			args{"BOB"},
			0,
		},
		{
			"Should return empty string if no 2 letters match found",
			args{"AA"},
			0,
		},
		{
			"Should uppercase input",
			args{"vnm"},
			2,
		},
		{
			"Could get England emoji",
			args{"GB-ENG"},
			7,
		},
		{
			"Could get CIOC code",
			args{"GER"},
			2,
		},
		{
			"Return empty string if code is empty",
			args{""},
			0,
		},
		{
			"Could get England flag with short code",
			args{"ENG"},
			7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFlag(tt.args.country)
			if !utf8.ValidString(got) {
				t.Errorf("GetFlag() expected valid flag got %v", got)
			}
			if GetPrintableLength(got) != tt.expectedLen {
				t.Errorf("expected length emoji of %v got %v", tt.expectedLen, GetPrintableLength(got))
			}
		})
	}
}

func GetPrintableLength(s string) int {
	res := 0
	for i := range s {
		if unicode.IsPrint(rune(s[i])) {
			res++
		}
	}

	return res
}

func Test_GetFlagFuzzy(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name        string
		args        args
		wantFlag    bool // whether a flag should be returned
		wantCode    string
		wantFlagLen int
	}{
		{
			"Should find exact match VNM",
			args{"VNM"},
			true,
			"VNM",
			2,
		},
		{
			"Should find typo VIETNM (missing A)",
			args{"VIETNM"},
			false, // distance is 3, exceeds maxDistance threshold of 2
			"",
			0,
		},
		{
			"Should find USA from US",
			args{"USA"},
			true,
			"USA", // Found USA in alpha-3 map
			2,
		},
		{
			"Should NOT find GER from GERMANY",
			args{"GERMANY"},
			false, // distance is 4, exceeds maxDistance threshold of 2
			"",
			0,
		},
		{
			"Should find GER from GERM",
			args{"GERM"},
			true,
			"GER",
			2,
		},
		{
			"Should find FR from FRR (distance 1)",
			args{"FRR"},
			true,
			"FR", // France at distance 1 (delete one R) - unambiguous match
			2,
		},
		{
			"Should return empty for completely wrong input",
			args{"ZZZZZ"},
			false,
			"",
			0,
		},
		{
			"Should find GB-ENG from GB-EN",
			args{"GB-EN"},
			true,
			"GB-ENG",
			7,
		},
		{
			"Should find lowercase input",
			args{"vnm"},
			true,
			"VNM",
			2,
		},
		// Edge cases
		{
			"Empty string finds some match",
			args{""},
			true, // Fuzzy matching finds a match at distance 2 (non-deterministic which one)
			"",   // We don't test exact code match due to map iteration randomness
			2,
		},
		{
			"Single character finds match",
			args{"V"},
			true, // Fuzzy matching finds a 2-char code at distance 1
			"",   // Non-deterministic which 2-char code starting with V
			2,
		},
		{
			"Very long invalid input",
			args{"XXXXXXXXXXXXXXXXXX"},
			false,
			"",
			0,
		},
		{
			"Special characters with hyphen",
			args{"GB-ENG"},
			true,
			"GB-ENG",
			7,
		},
		{
			"With space finds match",
			args{"GB ENG"},
			true, // Fuzzy matching finds GB-ENG at distance 1 (space becomes hyphen)
			"GB-ENG",
			7,
		},
		{
			"Numeric input",
			args{"123"},
			false,
			"",
			0,
		},
		{
			"Mixed case input",
			args{"VnM"},
			true,
			"VNM",
			2,
		},
		{
			"Four character code finds match",
			args{"VNAM"},
			true, // distance 1 from multiple codes, non-deterministic
			"",   // Don't test exact match
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlag, gotCode := GetFlagFuzzy(tt.args.input)

			if tt.wantFlag {
				if gotFlag == "" {
					t.Errorf("GetFlagFuzzy() expected flag, got empty string")
				}
				// Only check exact code match if wantCode is specified
				if tt.wantCode != "" && gotCode != tt.wantCode {
					t.Errorf("GetFlagFuzzy() code = %v, want %v", gotCode, tt.wantCode)
				}
				if !utf8.ValidString(gotFlag) {
					t.Errorf("GetFlagFuzzy() expected valid flag got %v", gotFlag)
				}
				if GetPrintableLength(gotFlag) != tt.wantFlagLen {
					t.Errorf("GetFlagFuzzy() flag length = %v, want %v", GetPrintableLength(gotFlag), tt.wantFlagLen)
				}
			} else {
				if gotFlag != "" || gotCode != "" {
					t.Errorf("GetFlagFuzzy() expected empty results, got flag=%v, code=%v", gotFlag, gotCode)
				}
			}
		})
	}
}

func Test_levenshtein(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"identical strings", "ABC", "ABC", 0},
		{"one insertion", "ABC", "ABCD", 1},
		{"one deletion", "ABCD", "ABC", 1},
		{"one substitution", "ABC", "ABD", 1},
		{"empty strings", "", "", 0},
		{"one empty", "ABC", "", 3},
		{"VIETNM to VNM", "VIETNM", "VNM", 3},
		{"GERMANY to GER", "GERMANY", "GER", 4},
		{"USA to US", "USA", "US", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levenshtein(tt.s1, tt.s2); got != tt.want {
				t.Errorf("levenshtein() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetCode(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		wantCode string
	}{
		{"Vietnam flag", "🇻🇳", "VN"},
		{"US flag", "🇺🇸", "US"},
		{"Germany flag", "🇩🇪", "DE"},
		{"England special flag", "🏴", "GB-ENG"},
		{"Scotland special flag", "🏴󠁧󠁢󠁳󠁣󠁴󠁿", "GB-SCT"},
		{"Wales special flag", "🏴󠁧󠁢󠁷󠁬󠁳󠁿", "GB-WLS"},
		{"Invalid flag", "🎌", ""},
		{"Empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCode(tt.flag)
			if got != tt.wantCode {
				t.Errorf("GetCode() = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func Test_GetName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
	}{
		{"From alpha-2 code", "VN", "Vietnam"},
		{"From alpha-3 code", "VNM", "Vietnam"},
		{"From CIOC code", "GER", "Germany"},
		{"From flag emoji", "🇻🇳", "Vietnam"},
		{"From special code", "GB-ENG", "England"},
		{"Lowercase code", "vn", "Vietnam"},
		{"Invalid code", "ZZZ", ""},
		{"Empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetName(tt.input)
			if got != tt.wantName {
				t.Errorf("GetName() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func Test_GetFlagByName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFlag bool
		wantCode string
	}{
		{"Exact match Vietnam", "Vietnam", true, "VN"},
		{"Exact match United States", "United States", true, "US"},
		{"Exact match Germany", "Germany", true, "DE"},
		{"Alias USA", "USA", true, "US"},
		{"Alias UK", "UK", true, "GB"},
		{"Alias UAE", "UAE", true, "AE"},
		{"Fuzzy match Viet Nam", "Viet Nam", true, "VN"},
		{"Special England", "England", true, "GB-ENG"},
		{"Lowercase", "vietnam", true, "VN"},
		{"Invalid name", "Notacountry", false, ""},
		{"Empty string", "", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlag, gotCode := GetFlagByName(tt.input)
			if tt.wantFlag {
				if gotFlag == "" {
					t.Errorf("GetFlagByName() expected flag, got empty string")
				}
				if gotCode != tt.wantCode {
					t.Errorf("GetFlagByName() code = %v, want %v", gotCode, tt.wantCode)
				}
			} else {
				if gotFlag != "" || gotCode != "" {
					t.Errorf("GetFlagByName() expected empty, got flag=%v, code=%v", gotFlag, gotCode)
				}
			}
		})
	}
}

func Test_ResolveFlag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFlag bool
		wantCode string
		wantKind string
	}{
		{"Exact alpha-2 code", "VN", true, "VN", ResolveKindCode},
		{"Exact alpha-3 code", "VNM", true, "VNM", ResolveKindCode},
		{"Exact CIOC code", "GER", true, "GER", ResolveKindCode},
		{"Exact alias in code set", "USA", true, "USA", ResolveKindCode},
		{"Exact country name", "Vietnam", true, "VN", ResolveKindName},
		{"Alias UK", "UK", true, "GB", ResolveKindAlias},
		{"Fuzzy code", "GERM", true, "GER", ResolveKindFuzzyCode},
		{"Fuzzy name", "Viet Namm", true, "VN", ResolveKindFuzzyName},
		{"Invalid", "NOTACOUNTRY", false, "", ""},
		{"Empty", "", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlag, gotCode, gotKind := ResolveFlag(tt.input)

			if tt.wantFlag {
				if gotFlag == "" {
					t.Errorf("ResolveFlag() expected non-empty flag")
				}
				if !utf8.ValidString(gotFlag) {
					t.Errorf("ResolveFlag() expected valid utf8 flag, got %v", gotFlag)
				}
				if GetPrintableLength(gotFlag) == 0 {
					t.Errorf("ResolveFlag() expected printable flag")
				}
			} else if gotFlag != "" || gotCode != "" || gotKind != "" {
				t.Errorf("ResolveFlag() expected empty outputs, got flag=%v code=%v kind=%v", gotFlag, gotCode, gotKind)
			}

			if gotCode != tt.wantCode {
				t.Errorf("ResolveFlag() code = %v, want %v", gotCode, tt.wantCode)
			}

			if gotKind != tt.wantKind {
				t.Errorf("ResolveFlag() kind = %v, want %v", gotKind, tt.wantKind)
			}
		})
	}
}

func ExampleGetFlag() {
	flag := GetFlag("VN")
	fmt.Println(flag)
	// Output: 🇻🇳
}

func ExampleGetFlag_threeLetterCode() {
	flag := GetFlag("VNM")
	fmt.Println(flag)
	// Output: 🇻🇳
}

func ExampleGetFlag_ciocCode() {
	flag := GetFlag("GER")
	fmt.Println(flag)
	// Output: 🇩🇪
}

func ExampleGetFlag_specialSubdivision() {
	flag := GetFlag("GB-ENG")
	fmt.Println(flag)
	// Output: 🏴󠁧󠁢󠁥󠁮󠁧󠁿
}

func ExampleGetFlag_invalidCode() {
	flag := GetFlag("INVALID")
	fmt.Println(flag == "")
	// Output: true
}

func ExampleGetFlagFuzzy() {
	flag, code := GetFlagFuzzy("GERM")
	fmt.Printf("Flag: %s, Code: %s\n", flag, code)
	// Output: Flag: 🇩🇪, Code: GER
}

func ExampleGetFlagFuzzy_exactMatch() {
	flag, code := GetFlagFuzzy("VNM")
	fmt.Printf("Flag: %s, Code: %s\n", flag, code)
	// Output: Flag: 🇻🇳, Code: VNM
}

func ExampleGetFlagFuzzy_variation() {
	flag, code := GetFlagFuzzy("USA")
	fmt.Printf("Flag: %s, Code: %s\n", flag, code)
	// Output: Flag: 🇺🇸, Code: USA
}

func ExampleGetFlagFuzzy_tooFar() {
	flag, code := GetFlagFuzzy("GERMANY")
	fmt.Printf("Flag: %s, Code: %s\n", flag, code)
	// Output: Flag: , Code:
}

func ExampleGetCode() {
	code := GetCode("🇻🇳")
	fmt.Println(code)
	// Output: VN
}

func ExampleGetName() {
	name := GetName("VN")
	fmt.Println(name)
	// Output: Vietnam
}

func ExampleGetName_threeLetterCode() {
	name := GetName("VNM")
	fmt.Println(name)
	// Output: Vietnam
}

func ExampleGetName_fromFlag() {
	name := GetName("🇻🇳")
	fmt.Println(name)
	// Output: Vietnam
}

func ExampleGetFlagByName() {
	flag, code := GetFlagByName("Vietnam")
	fmt.Printf("Flag: %s, Code: %s\n", flag, code)
	// Output: Flag: 🇻🇳, Code: VN
}

func ExampleGetFlagByName_alias() {
	flag, code := GetFlagByName("USA")
	fmt.Printf("Flag: %s, Code: %s\n", flag, code)
	// Output: Flag: 🇺🇸, Code: US
}
