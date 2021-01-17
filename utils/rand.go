package utils

import (
	"math"
	"math/rand"
	"time"

	"github.com/jrapoport/sillyname-go"
	"github.com/lucasb-eyer/go-colorful"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomUsername generates a random username
func RandomUsername(maxLength int) string {
	return sillyname.GenerateUsernameN(maxLength)
}

// RandomColor generates a random hex color
func RandomColor() string {
	i := rand.Float64()
	c := colorful.Hsv(math.Mod(i*0.618033988749895, 1.0),
		0.5,
		math.Sqrt(1.0-math.Mod(i*0.618033988749895, 0.5)))
	return c.Hex()
}

// RandomPIN generates a random numerical PIN code of length
func RandomPIN(length int) string {
	const pool = "1234567890"
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes)
}
