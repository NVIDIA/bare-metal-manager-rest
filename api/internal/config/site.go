package config

// SiteConfig holds configuration for components and services running from site
type SiteConfig struct {
	Disconnected bool
	PhoneHomeUrl string
}

// NewSiteConfig initializes and returns a configuration object for site
func NewSiteConfig(disconnected bool, phoneHomeUrl string) *SiteConfig {
	return &SiteConfig{
		Disconnected: disconnected,
		PhoneHomeUrl: phoneHomeUrl,
	}
}
