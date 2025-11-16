package argon2

type Argon2 struct{}

type argonConfig struct {
	SaltLength uint32
	KeyLength  uint32
	Iterations uint32
	Memory     uint32
	Parallel   uint8
}
