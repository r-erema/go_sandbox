package tls

type (
	messageType               [1]byte
	cipherSuite               [2]byte
	supportedKeyExchangeGroup [2]byte
	signatureAlgorithm        [2]byte

	supportedTLSVersion [2]byte

	pskKeyExchangeMode [1]byte

	publicKey struct {
		payload       [32]byte
		exchangeGroup supportedKeyExchangeGroup
	}

	ClientApplicationKey []byte
	ClientApplicationIV  []byte
	ServerApplicationKey []byte
	ServerApplicationIV  []byte
)

func clientType() messageType {
	return messageType{0x01}
}

func serverType() messageType {
	return messageType{0x02}
}

func tlsAes256GcmSha384() cipherSuite {
	return cipherSuite{0x13, 0x02}
}

func x25519() supportedKeyExchangeGroup {
	return supportedKeyExchangeGroup{0x00, 0x1d}
}

func rsaPssRsaeSha512() signatureAlgorithm {
	return signatureAlgorithm{0x08, 0x06}
}

func tls13() supportedTLSVersion {
	return supportedTLSVersion{0x03, 0x04}
}

func pskWithECDHEKeyEstablishment() pskKeyExchangeMode {
	return pskKeyExchangeMode{0x01}
}
