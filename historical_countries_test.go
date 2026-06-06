package emojiflags

import (
	"testing"
)

func Test_GetFlag_Historical(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		// USSR
		{"USSR alpha-2", "SU", "\U0001F1F8\U0001F1FA"},
		{"USSR alpha-3", "SUN", "\U0001F1F8\U0001F1FA"},
		{"USSR CIOC", "URS", "\U0001F1F8\U0001F1FA"},
		{"USSR alpha-4", "SUHH", "\U0001F1F8\U0001F1FA"},
		{"USSR lowercase alpha-2", "su", "\U0001F1F8\U0001F1FA"},

		// Yugoslavia
		{"Yugoslavia alpha-2", "YU", "\U0001F1FE\U0001F1FA"},
		{"Yugoslavia alpha-3", "YUG", "\U0001F1FE\U0001F1FA"},
		{"Yugoslavia CIOC", "YUG", "\U0001F1FE\U0001F1FA"},
		{"Yugoslavia alpha-4", "YUCS", "\U0001F1FE\U0001F1FA"},

		// Czechoslovakia
		{"Czechoslovakia ambiguous alpha-2 defaults to first user", "CS", "\U0001F1E8\U0001F1F8"},
		{"Czechoslovakia alpha-3", "CSK", "\U0001F1E8\U0001F1F8"},
		{"Czechoslovakia CIOC", "TCH", "\U0001F1E8\U0001F1F8"},
		{"Czechoslovakia alpha-4", "CSHH", "\U0001F1E8\U0001F1F8"},

		// Serbia and Montenegro
		{"Serbia and Montenegro alpha-3", "SCG", "\U0001F1E8\U0001F1F8"},
		{"Serbia and Montenegro CIOC", "SCG", "\U0001F1E8\U0001F1F8"},
		{"Serbia and Montenegro alpha-4", "CSXX", "\U0001F1E8\U0001F1F8"},

		// East Germany
		{"East Germany alpha-2", "DD", "\U0001F1E9\U0001F1E9"},
		{"East Germany alpha-3", "DDR", "\U0001F1E9\U0001F1E9"},
		{"East Germany CIOC", "GDR", "\U0001F1E9\U0001F1E9"},
		{"East Germany alpha-4", "DDDE", "\U0001F1E9\U0001F1E9"},

		// Zaire
		{"Zaire alpha-2", "ZR", "\U0001F1FF\U0001F1F7"},
		{"Zaire alpha-3", "ZAR", "\U0001F1FF\U0001F1F7"},
		{"Zaire CIOC", "ZAI", "\U0001F1FF\U0001F1F7"},
		{"Zaire alpha-4", "ZRCD", "\U0001F1FF\U0001F1F7"},

		// North/South Vietnam intentionally not supported.
		{"North Vietnam absent", "VD", ""},
		{"South Vietnam absent (current VN)", "VN", "\U0001F1FB\U0001F1F3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFlag(tt.code)
			if got != tt.expected {
				t.Errorf("GetFlag(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

func Test_GetName_Historical(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"USSR alpha-2", "SU", "Soviet Union"},
		{"USSR alpha-3", "SUN", "Soviet Union"},
		{"USSR CIOC", "URS", "Soviet Union"},
		{"USSR alpha-4", "SUHH", "Soviet Union"},
		{"USSR lowercase", "su", "Soviet Union"},

		{"Yugoslavia alpha-2", "YU", "Yugoslavia"},
		{"Yugoslavia alpha-3", "YUG", "Yugoslavia"},
		{"Yugoslavia alpha-4", "YUCS", "Yugoslavia"},

		// CS disambiguation: bare CS defaults to Czechoslovakia;
		// SCG and CSXX must resolve to Serbia and Montenegro.
		{"CS defaults to Czechoslovakia", "CS", "Czechoslovakia"},
		{"CSK resolves to Czechoslovakia", "CSK", "Czechoslovakia"},
		{"TCH resolves to Czechoslovakia", "TCH", "Czechoslovakia"},
		{"CSHH resolves to Czechoslovakia", "CSHH", "Czechoslovakia"},
		{"SCG resolves to Serbia and Montenegro", "SCG", "Serbia and Montenegro"},
		{"CSXX resolves to Serbia and Montenegro", "CSXX", "Serbia and Montenegro"},

		{"East Germany alpha-2", "DD", "East Germany"},
		{"East Germany alpha-3", "DDR", "East Germany"},
		{"East Germany CIOC", "GDR", "East Germany"},
		{"East Germany alpha-4", "DDDE", "East Germany"},

		{"Zaire alpha-2", "ZR", "Zaire"},
		{"Zaire alpha-3", "ZAR", "Zaire"},
		{"Zaire CIOC", "ZAI", "Zaire"},
		{"Zaire alpha-4", "ZRCD", "Zaire"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetName(tt.input)
			if got != tt.expected {
				t.Errorf("GetName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func Test_GetCountryInfo_Historical(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantAlpha2 string
		wantAlpha3 string
		wantCioc   string
		wantName   string
	}{
		{"USSR via alpha-2", "SU", "SU", "SUN", "URS", "Soviet Union"},
		{"USSR via alpha-3", "SUN", "SU", "SUN", "URS", "Soviet Union"},
		{"USSR via CIOC", "URS", "SU", "SUN", "URS", "Soviet Union"},
		{"USSR via alpha-4", "SUHH", "SU", "SUN", "URS", "Soviet Union"},

		{"Yugoslavia via alpha-2", "YU", "YU", "YUG", "YUG", "Yugoslavia"},
		{"Yugoslavia via alpha-4", "YUCS", "YU", "YUG", "YUG", "Yugoslavia"},

		{"CS via alpha-2 defaults to Czechoslovakia", "CS", "CS", "CSK", "TCH", "Czechoslovakia"},
		{"Czechoslovakia via CSK", "CSK", "CS", "CSK", "TCH", "Czechoslovakia"},
		{"Czechoslovakia via TCH", "TCH", "CS", "CSK", "TCH", "Czechoslovakia"},
		{"Czechoslovakia via CSHH", "CSHH", "CS", "CSK", "TCH", "Czechoslovakia"},
		{"Serbia and Montenegro via SCG", "SCG", "CS", "SCG", "SCG", "Serbia and Montenegro"},
		{"Serbia and Montenegro via CSXX", "CSXX", "CS", "SCG", "SCG", "Serbia and Montenegro"},

		{"East Germany via DD", "DD", "DD", "DDR", "GDR", "East Germany"},
		{"East Germany via GDR", "GDR", "DD", "DDR", "GDR", "East Germany"},
		{"East Germany via DDDE", "DDDE", "DD", "DDR", "GDR", "East Germany"},

		{"Zaire via ZR", "ZR", "ZR", "ZAR", "ZAI", "Zaire"},
		{"Zaire via ZAI", "ZAI", "ZR", "ZAR", "ZAI", "Zaire"},
		{"Zaire via ZRCD", "ZRCD", "ZR", "ZAR", "ZAI", "Zaire"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCountryInfo(tt.input)
			if got.Alpha2 != tt.wantAlpha2 {
				t.Errorf("Alpha2 = %q, want %q", got.Alpha2, tt.wantAlpha2)
			}
			if got.Alpha3 != tt.wantAlpha3 {
				t.Errorf("Alpha3 = %q, want %q", got.Alpha3, tt.wantAlpha3)
			}
			if got.CIOC != tt.wantCioc {
				t.Errorf("CIOC = %q, want %q", got.CIOC, tt.wantCioc)
			}
			if got.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", got.Name, tt.wantName)
			}
		})
	}
}

func Test_LookupHistorical_Unknown(t *testing.T) {
	// Codes that are not in the historical list must return an empty record.
	for _, code := range []string{"ZZ", "ZZZ", "ZZZZ", "ZZZZZ", "ABCD", "", "FR"} {
		if _, ok := lookupHistorical(code); ok {
			t.Errorf("lookupHistorical(%q) should return false", code)
		}
	}
}

func Test_HistoricalCountryRecordsAreConsistent(t *testing.T) {
	// Every entry must have non-empty Alpha2 (2), Alpha3 (3), Alpha4 (4),
	// Cioc (3), and Name fields, and the alpha-2/3 codes must be
	// uppercase ASCII letters of the correct length.
	for _, hc := range HistoricalCountries {
		if len(hc.Alpha2) != 2 {
			t.Errorf("%s: Alpha2 %q must be 2 chars", hc.Name, hc.Alpha2)
		}
		if len(hc.Alpha3) != 3 {
			t.Errorf("%s: Alpha3 %q must be 3 chars", hc.Name, hc.Alpha3)
		}
		if len(hc.Alpha4) != 4 {
			t.Errorf("%s: Alpha4 %q must be 4 chars", hc.Name, hc.Alpha4)
		}
		if len(hc.Cioc) != 3 {
			t.Errorf("%s: Cioc %q must be 3 chars", hc.Name, hc.Cioc)
		}
		if hc.Name == "" {
			t.Errorf("%s: Name must not be empty", hc.Alpha4)
		}
	}
}

func Test_HistoricalCountryAlpha4sAreUnique(t *testing.T) {
	// ISO 3166-3 alpha-4 codes are unique by definition; guard against
	// accidental duplicates in the source-of-truth slice.
	seen := make(map[string]string, len(HistoricalCountries))
	for _, hc := range HistoricalCountries {
		if prev, ok := seen[hc.Alpha4]; ok {
			t.Errorf("duplicate Alpha4 %q used by %q and %q", hc.Alpha4, prev, hc.Name)
		}
		seen[hc.Alpha4] = hc.Name
	}
}
