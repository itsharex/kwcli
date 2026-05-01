package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ToMap converts KWDBConfig to a map[string]interface{}
func (c *KWDBConfig) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"sql_port":   c.SQLPort,
		"http_port":  c.HTTPPort,
		"data_dir":   c.DataDir,
		"log_dir":    c.LogDir,
		"insecure":   c.Insecure,
		"listen_addr": c.ListenAddr,
		"http_addr":  c.HTTPAddr,
	}
}

// KWDB default configuration
var DefaultKWDBConfig = KWDBConfig{
	SQLPort:   26257,
	HTTPPort:  8080,
	Insecure:  true,
	ListenAddr: "0.0.0.0:26257",
	HTTPAddr:  "0.0.0.0:8080",
}

// KWDBConfig represents KWDB configuration
type KWDBConfig struct {
	SQLPort    int    `mapstructure:"sql_port"`
	HTTPPort   int    `mapstructure:"http_port"`
	DataDir    string `mapstructure:"data_dir"`
	LogDir     string `mapstructure:"log_dir"`
	Insecure   bool   `mapstructure:"insecure"`
	ListenAddr string `mapstructure:"listen_addr"`
	HTTPAddr   string `mapstructure:"http_addr"`
}

// GetHomeDir returns the KWCLI home directory
func GetHomeDir() string {
	// Check viper config first
	if viper.IsSet("home") {
		return viper.GetString("home")
	}
	// Check environment variable
	if home := os.Getenv("KWCLI_HOME"); home != "" {
		return home
	}
	// Default to ~/.kwcli
	home, err := os.UserHomeDir()
	if err != nil {
		return ".kwcli"
	}
	return filepath.Join(home, ".kwcli")
}

// GetComponentsDir returns the components directory
func GetComponentsDir() string {
	return filepath.Join(GetHomeDir(), "components")
}

// GetDataDir returns the data directory
func GetDataDir() string {
	return filepath.Join(GetHomeDir(), "data")
}

// GetBinDir returns the bin directory
func GetBinDir() string {
	return filepath.Join(GetHomeDir(), "bin")
}

// GetKWDBConfigPath returns the path to KWDB config file
func GetKWDBConfigPath() string {
	return filepath.Join(GetComponentsDir(), "kwdb", "config.yaml")
}

// LoadKWDBConfig loads KWDB configuration
func LoadKWDBConfig() (*KWDBConfig, error) {
	configPath := GetKWDBConfigPath()
	
	v := viper.New()
	v.SetConfigFile(configPath)
	
	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return &DefaultKWDBConfig, nil
		}
		return nil, err
	}
	
	var config KWDBConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	// Apply defaults
	if config.SQLPort == 0 {
		config.SQLPort = DefaultKWDBConfig.SQLPort
	}
	if config.HTTPPort == 0 {
		config.HTTPPort = DefaultKWDBConfig.HTTPPort
	}
	if config.ListenAddr == "" {
		config.ListenAddr = DefaultKWDBConfig.ListenAddr
	}
	if config.HTTPAddr == "" {
		config.HTTPAddr = DefaultKWDBConfig.HTTPAddr
	}
	
	return &config, nil
}

// SaveKWDBConfig saves KWDB configuration
func SaveKWDBConfig(config *KWDBConfig) error {
	configPath := GetKWDBConfigPath()
	
	// Ensure directory exists
	os.MkdirAll(filepath.Dir(configPath), 0755)
	
	v := viper.New()
	v.SetConfigFile(configPath)
	v.Set("sql_port", config.SQLPort)
	v.Set("http_port", config.HTTPPort)
	v.Set("data_dir", config.DataDir)
	v.Set("log_dir", config.LogDir)
	v.Set("insecure", config.Insecure)
	v.Set("listen_addr", config.ListenAddr)
	v.Set("http_addr", config.HTTPAddr)
	
	return v.WriteConfig()
}