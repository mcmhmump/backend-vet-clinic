package domain

type AppConfig struct {
	IPAccess IPAccessConfig `mapstructure:"ip_access" yaml:"ip_access"`
}

type IPAccessConfig struct {
	DefaultPolicy string   `mapstructure:"default_policy" yaml:"default_policy"`
	Blacklist     []string `mapstructure:"blacklist" yaml:"blacklist"`
	Whitelist     []string `mapstructure:"whitelist" yaml:"whitelist"`
	Greylist      []string `mapstructure:"greylist" yaml:"greylist"`
}
