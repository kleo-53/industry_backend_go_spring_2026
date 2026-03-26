package main

func reverseRunes(s string) string {
	lenS := len([]rune(s))
	res := make([]rune, lenS)
	for i, runa := range []rune(s) {
		res[lenS - 1 - i] = runa
	}
	return string(res)
}
