package diffator

// longestCommonSubstr finds the longest common substring between two strings, properly handling Unicode.
func longestCommonSubstr(s1, s2 string) (string, int) {
	var maxLength, endIndex int

	runes1, runes2 := []rune(s1), []rune(s2)
	lenR1, lenR2 := len(runes1), len(runes2)

	// Create a 2D slice to store lengths of longest common suffixes
	lcsSuffix := make([][]int, lenR1+1)
	for i := range lcsSuffix {
		lcsSuffix[i] = make([]int, lenR2+1)
	}

	// Build lcsSuffix in bottom up manner
	for i := 0; i <= lenR1; i++ {
		for j := 0; j <= lenR2; j++ {
			if i == 0 || j == 0 {
				lcsSuffix[i][j] = 0
			} else if runes1[i-1] == runes2[j-1] {
				lcsSuffix[i][j] = lcsSuffix[i-1][j-1] + 1
				if maxLength < lcsSuffix[i][j] {
					maxLength = lcsSuffix[i][j]
					endIndex = i
				}
			} else {
				lcsSuffix[i][j] = 0
			}
		}
	}

	// If no common substring exists
	if maxLength == 0 {
		return "", 0
	}

	// Return the longest common substring
	return string(runes1[endIndex-maxLength : endIndex]), endIndex - maxLength
}
