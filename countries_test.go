package emojiflags

import (
	"fmt"
	"testing"
	"unicode"
	"unicode/utf8"
)

func Test_GetFlag(t *testing.T) {
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

func Test_GetCountryInfo(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantAlpha2 string
		wantAlpha3 string
		wantCioc   string
		wantName   string
	}{
		// Exact matches
		{"From alpha-2 code", "VN", "VN", "VNM", "VIE", "Vietnam"},
		{"From alpha-3 code", "VNM", "VN", "VNM", "VIE", "Vietnam"},
		{"From CIOC code", "GER", "DE", "DEU", "GER", "Germany"},
		{"From country name", "Vietnam", "VN", "VNM", "VIE", "Vietnam"},
		{"From alias", "UK", "GB", "GBR", "GBR", "United Kingdom"},
		{"From flag emoji", "🇺🇸", "US", "USA", "USA", "United States"},
		{"Lowercase input", "vietnam", "VN", "VNM", "VIE", "Vietnam"},
		{"Special subdivision", "GB-ENG", "GB-ENG", "", "", "England"},

		// Fuzzy matches
		{"Fuzzy alpha-3", "GERM", "DE", "DEU", "GER", "Germany"},
		{"Fuzzy alpha-2", "VNN", "VN", "VNM", "VIE", "Vietnam"},
		{"Fuzzy country name", "Germany", "DE", "DEU", "GER", "Germany"},
		{"Fuzzy alias", "Englnd", "GB-ENG", "", "", "England"},

		// No match
		{"Invalid input", "NOTACOUNTRY", "", "", "", ""},
		{"Empty string", "", "", "", "", ""},
		{"Too fuzzy", "XYZXYZXYZ", "", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCountryInfo(tt.input)
			if got.Alpha2 != tt.wantAlpha2 {
				t.Errorf("GetCountryInfo().Alpha2 = %v, want %v", got.Alpha2, tt.wantAlpha2)
			}
			if got.Alpha3 != tt.wantAlpha3 {
				t.Errorf("GetCountryInfo().Alpha3 = %v, want %v", got.Alpha3, tt.wantAlpha3)
			}
			if got.CIOC != tt.wantCioc {
				t.Errorf("GetCountryInfo().CIOC = %v, want %v", got.CIOC, tt.wantCioc)
			}
			if got.Name != tt.wantName {
				t.Errorf("GetCountryInfo().Name = %v, want %v", got.Name, tt.wantName)
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

func ExampleGetCountryInfo() {
	info := GetCountryInfo("Vietnam")
	fmt.Printf("Alpha2: %s, Alpha3: %s, CIOC: %s, Name: %s\n", info.Alpha2, info.Alpha3, info.CIOC, info.Name)
	// Output: Alpha2: VN, Alpha3: VNM, CIOC: VIE, Name: Vietnam
}

func ExampleGetCountryInfo_fromCode() {
	info := GetCountryInfo("VNM")
	fmt.Printf("Alpha2: %s, Alpha3: %s, CIOC: %s, Name: %s\n", info.Alpha2, info.Alpha3, info.CIOC, info.Name)
	// Output: Alpha2: VN, Alpha3: VNM, CIOC: VIE, Name: Vietnam
}

func ExampleGetCountryInfo_fromFlag() {
	info := GetCountryInfo("🇺🇸")
	fmt.Printf("Alpha2: %s, Alpha3: %s, CIOC: %s, Name: %s\n", info.Alpha2, info.Alpha3, info.CIOC, info.Name)
	// Output: Alpha2: US, Alpha3: USA, CIOC: USA, Name: United States
}

func ExampleGetCountryInfo_fuzzy() {
	info := GetCountryInfo("GERM")
	fmt.Printf("Alpha2: %s, Name: %s\n", info.Alpha2, info.Name)
	// Output: Alpha2: DE, Name: Germany
}
