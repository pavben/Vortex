package aesstream

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// AesStream is a wrapper for a ReadWriter with two-way AES encryption.
type AesStream struct {
	readWriter   io.ReadWriter
	cfbEncrypter cipher.Stream
	cfbDecrypter cipher.Stream
}

// NewAesStream creates a new AesStream.
func NewAesStream(readWriter io.ReadWriter, aesKey, iv []byte) (*AesStream, error) {
	if len(iv) != aes.BlockSize {
		panic("Length of iv in NewAesCrypter does not match the AES block size")
	}
	// Create the AES cipher
	aesCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	// TODO: Is it safe to initialize both the encrypter and the decrypter to the same IV?
	aesStream := &AesStream{
		readWriter:   readWriter,
		cfbEncrypter: cipher.NewCFBEncrypter(aesCipher, iv),
		cfbDecrypter: cipher.NewCFBDecrypter(aesCipher, iv),
	}
	return aesStream, nil
}

// Read will read as many bytes as possible from the ReadWriter into p and decrypt them, returning the number of bytes read.
func (as *AesStream) Read(p []byte) (int, error) {
	n, err := as.readWriter.Read(p)
	if err != nil {
		return 0, err
	}
	as.cfbDecrypter.XORKeyStream(p[:n], p[:n])
	return n, nil
}

// Write will encrypt the bytes in p and write them to the ReadWriter.
func (as *AesStream) Write(p []byte) (int, error) {
	pEncrypted := make([]byte, len(p))
	as.cfbEncrypter.XORKeyStream(pEncrypted, p)
	n, err := as.readWriter.Write(pEncrypted)
	if err != nil {
		return 0, err
	}
	return n, nil
}
