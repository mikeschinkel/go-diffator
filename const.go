package diffator

// MinSubstrLen is the default length of a "common substring" as defined by being
// longer than 3 characters and not being a space.
const MinSubstrLen = 3

// LeftRightFormat is default format used to format the differences found by
// `CompareStrings()` inline within the output string, would be the formatted when comparing these two (2)
// strings with the default format below:
//
//	Left String:    "this shows left content inline"
//	Right String:   "this shows right content inline"
//	Compare Output: "this shows (left/right) content inline"
const LeftRightFormat = "<(%s/%s)>"
