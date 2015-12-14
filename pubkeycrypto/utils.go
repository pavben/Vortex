package pubkeycrypto

import "crypto"

var sha512HashFunc = crypto.SHA512

func sha512Sum(data []byte) []byte {
	hash := sha512HashFunc.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func sha1Sum(data []byte) []byte {
	hash := crypto.SHA1.New()
	hash.Write(data)
	return hash.Sum(nil)
}
