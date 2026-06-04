package emojiflags

import "testing"

func BenchmarkGetFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFlag("VNM")
	}
}

func BenchmarkGetFlag2Letter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFlag("VN")
	}
}

func BenchmarkGetFlagCIOC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFlag("GER")
	}
}

func BenchmarkGetFlagSpecial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFlag("GB-ENG")
	}
}

func BenchmarkGetFlagInvalid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFlag("INVALID")
	}
}

func BenchmarkGetCountryInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCountryInfo("VNM")
	}
}

func BenchmarkGetCountryInfoByName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCountryInfo("Vietnam")
	}
}

func BenchmarkGetCountryInfoByFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCountryInfo("🇻🇳")
	}
}

func BenchmarkGetCountryInfoFuzzy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCountryInfo("GERM")
	}
}

func BenchmarkGetCountryInfoInvalid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCountryInfo("INVALID")
	}
}

func BenchmarkGetName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetName("VN")
	}
}

func BenchmarkGetNameByAlpha3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetName("VNM")
	}
}

func BenchmarkGetNameByFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetName("🇻🇳")
	}
}

func BenchmarkLevenshtein(b *testing.B) {
	for i := 0; i < b.N; i++ {
		levenshtein("VIETNM", "VNM")
	}
}

func BenchmarkLevenshteinShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		levenshtein("USA", "US")
	}
}

func BenchmarkLevenshteinIdentical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		levenshtein("VNM", "VNM")
	}
}
