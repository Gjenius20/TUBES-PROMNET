package autograder

import (
	"math"
	"strconv"
	"strings"
)

// floatTolerance adalah toleransi relatif untuk membandingkan token numerik.
const floatTolerance = 1e-6

// Normalize menyeragamkan newline (CRLF -> LF), membuang spasi di akhir setiap baris
// dan baris kosong di akhir output. Dipakai untuk menampilkan output ke pengguna.
func Normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

// Match membandingkan output program dengan output yang diharapkan.
//
// Aturan (tolerance matcher):
//  1. Perbedaan whitespace/newline diabaikan; output dipecah menjadi token.
//  2. Urutan token harus sama persis (sequence matching) dan jumlahnya harus sama.
//  3. Token numerik dianggap sama bila selisihnya dalam toleransi relatif 1e-6.
//  4. Token non-numerik (keyword/teks) dibandingkan persis, case-sensitive.
func Match(actual, expected string) bool {
	a := strings.Fields(actual)
	e := strings.Fields(expected)
	if len(a) != len(e) {
		return false
	}
	for i := range e {
		if !tokenEqual(a[i], e[i]) {
			return false
		}
	}
	return true
}

func tokenEqual(actual, expected string) bool {
	if actual == expected {
		return true
	}
	af, aErr := strconv.ParseFloat(actual, 64)
	ef, eErr := strconv.ParseFloat(expected, 64)
	if aErr != nil || eErr != nil {
		return false
	}
	if math.IsNaN(af) || math.IsNaN(ef) || math.IsInf(af, 0) || math.IsInf(ef, 0) {
		return false
	}
	return math.Abs(af-ef) <= floatTolerance*math.Max(1, math.Abs(ef))
}
