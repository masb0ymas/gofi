package constant

func AllowedOrigin() []string {
	local := []string{"http://localhost:3000", "http://localhost:3333"}
	internal := []string{"https://example.com"}

	result := append(local, internal...)

	return result
}
