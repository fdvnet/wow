package pow

import (
	"crypto/sha256"
	"errors"
	"strconv"
	"time"
)

const timeout = time.Minute

func isValid(zeroCount int, b []byte) bool {
	i := 0
	for zeroCount >= 8 {
		if b[i] != 0 {
			return false
		}
		i++
		zeroCount -= 8
	}

	var mask byte
	switch zeroCount {
	case 0:
		mask = 0x00
	case 1:
		mask = 0x80
	case 2:
		mask = 0xC0
	case 3:
		mask = 0xE0
	case 4:
		mask = 0xF0
	case 5:
		mask = 0xF8
	case 6:
		mask = 0xFC
	case 7:
		mask = 0xFE
	}
	return (b[i] & mask) == 0
}

func IsValid(nonce, solution []byte, zeroCount int) bool {
	if zeroCount > len(solution)*8 {
		return false
	}
	hasher := sha256.New()
	hasher.Write(solution)
	hasher.Write(nonce)
	hash := hasher.Sum(nil)
	return isValid(zeroCount, hash)
}

var errTimeOver = errors.New("time for calculation is over")

// Calculate
// it is possible to make parallel calculation using goroutines
func Calculate(nonce []byte, zeroCount int) ([]byte, error) {
	var sol []byte
	hasher := sha256.New()
	tc := time.NewTimer(timeout)
	var i int
	for {
		i++
		select {
		case <-tc.C:
			return nil, errTimeOver
		default:
			hasher.Reset()
			sol = []byte(strconv.Itoa(i))
			hasher.Write(sol)
			hasher.Write(nonce)
			hash := hasher.Sum(nil)
			if isValid(zeroCount, hash) {
				return sol, nil
			}
		}
	}
}
