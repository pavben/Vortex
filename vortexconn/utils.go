package vortexconn

import "crypto/rand"

func generateRandomBytes(numBytes int) ([]byte, error) {
	b := make([]byte, numBytes)
	_, err := rand.Reader.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
