package pubkeycrypto

// KeyPair is a holder for a private key and a public key.
type KeyPair struct {
	PrivateKey *PrivateKey
	PublicKey  *PublicKey
}

// GenerateKeyPair generates a private-public keypair.
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.GetPublicKey()
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
