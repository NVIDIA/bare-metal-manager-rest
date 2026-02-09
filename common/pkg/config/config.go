package config

// PayloadEncryptionConfig holds configuration for payload encryption
type PayloadEncryptionConfig struct {
	EncryptionKey string
}

// NewPayloadEncryptionConfig initializes and returns a configuration object for payload encryption
func NewPayloadEncryptionConfig(encryptionKey string) *PayloadEncryptionConfig {
	return &PayloadEncryptionConfig{
		EncryptionKey: encryptionKey,
	}
}
