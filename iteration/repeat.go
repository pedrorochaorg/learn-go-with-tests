package iteration

// Repeat concatenates 'repatCount' the same value of 'character' in a single string
func Repeat(character string, repeatCount int) string {
	var repeated string
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return repeated
}
