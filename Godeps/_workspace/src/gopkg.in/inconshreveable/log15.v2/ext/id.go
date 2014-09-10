package ext

import (
	"fmt"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandId creates a random identifier of the requested length.
// Useful for assigning mostly-unique identifiers for logging
// and identification that are unlikely to collide because of
// short lifespan or low set cardinality
func RandId(idlen int) string {
	b := make([]byte, idlen)
	var randVal uint32
	for i := 0; i < idlen; i++ {
		byteIdx := i % 4
		if byteIdx == 0 {
			randVal = r.Uint32()
		}
		b[i] = byte((randVal >> (8 * uint(byteIdx))) & 0xFF)
	}
	return fmt.Sprintf("%x", b)
}
