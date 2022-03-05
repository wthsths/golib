package strings

func MaskRange(s string, maskStart, maskEnd int, maskChar rune) string {
	rs := []rune(s)

	for i := 6; i < len(rs); i++ {
		if i >= maskStart && i <= maskEnd {
			rs[i] = maskChar
		}
	}
	return string(rs)
}
