package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultKWDBConfig(t *testing.T) {
	cfg := DefaultKWDBConfig

	if cfg.SQLPort != 26257 {
		t.Errorf("Expected SQLPort 26257, got %d", cfg.SQLPort)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("Expected HTTPPort 8080, got %d", cfg.HTTPPort)
	}
	if cfg.Insecure != true {
		t.Errorf("Expected Insecure true, got %v", cfg.Insecure)
	}
	if cfg.ListenAddr != "0.0.0.0:26257" {
		t.Errorf("Expected ListenAddr '0.0.0.0:26257', got '%s'", cfg.ListenAddr)
	}
	if cfg.HTTPAddr != "0.0.0.0:8080" {
		t.Errorf("Expected HTTPAddr '0.0.0.0:8080', got '%s'", cfg.HTTPAddr)
	}
}

func TestGetHomeDir(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantHome string
	}{
		{"default", "", ""},
		{"custom", "/custom/path", "/custom/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("KWCLI_HOME", tt.envValue)
				defer os.Unsetenv("KWCLI_HOME")
			} else {
				os.Unsetenv("KWCLI_HOME")
			}

			home := GetHomeDir()
			if tt.wantHome != "" && home != tt.wantHome {
				t.Errorf("Expected home '%s', got '%s'", tt.wantHome, home)
			}
			// For default case, just check it's not empty
			if tt.wantHome == "" && home == "" {
				t.Error("Expected home to be set by default")
			}
		})
	}
}

func TestGetComponentsDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	componentsDir := GetComponentsDir()
	expected := filepath.Join(testHome, "components")

	if componentsDir != expected {
		t.Errorf("Expected '%s', got '%s'", expected, componentsDir)
	}
}

func TestGetDataDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	dataDir := GetDataDir()
	expected := filepath.Join(testHome, "data")

	if dataDir != expected {
		t.Errorf("Expected '%s', got '%s'", expected, dataDir)
	}
}

func TestGetBinDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	binDir := GetBinDir()
	expected := filepath.Join(testHome, "bin")

	if binDir != expected {
		t.Errorf("Expected '%s', got '%s'", expected, binDir)
	}
}

func TestGetKWDBConfigPath(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	configPath := GetKWDBConfigPath()
	expected := filepath.Join(testHome, "components", "kwdb", "config.yaml")

	if configPath != expected {
		t.Errorf("Expected '%s', got '%s'", expected, configPath)
	}
}

func TestSaveAndLoadKWDBConfig(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	// Create directory
	os.MkdirAll(filepath.Join(testHome, "components", "kwdb"), 0755)

	// Test config
	testCfg := KWDBConfig{
		SQLPort:    26257,
		HTTPPort:   8080,
		DataDir:    "/data/kwdb",
		LogDir:     "/logs/kwdb",
		Insecure:   true,
		ListenAddr: "0.0.0.0:26257",
		HTTPAddr:   "0.0.0.0:8080",
	}

	// Save config
	if err := SaveKWDBConfig(&testCfg); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := LoadKWDBConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify
	if loadedCfg.SQLPort != testCfg.SQLPort {
		t.Errorf("Expected SQLPort %d, got %d", testCfg.SQLPort, loadedCfg.SQLPort)
	}
	if loadedCfg.HTTPPort != testCfg.HTTPPort {
		t.Errorf("Expected HTTPPort %d, got %d", testCfg.HTTPPort, loadedCfg.HTTPPort)
	}
	if loadedCfg.DataDir != testCfg.DataDir {
		t.Errorf("Expected DataDir '%s', got '%s'", testCfg.DataDir, loadedCfg.DataDir)
	}
	if loadedCfg.LogDir != testCfg.LogDir {
		t.Errorf("Expected LogDir '%s', got '%s'", testCfg.LogDir, loadedCfg.LogDir)
	}
	if loadedCfg.Insecure != testCfg.Insecure {
		t.Errorf("Expected Insecure %v, got %v", testCfg.Insecure, loadedCfg.Insecure)
	}
}

func TestLoadKWDBConfig_NotExists(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	// Don't create config file - should return default
	cfg, err := LoadKWDBConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return default values
	if cfg.SQLPort != 26257 {
		t.Errorf("Expected default SQLPort 26257, got %d", cfg.SQLPort)
	}
}

func TestKWDBConfigFields(t *testing.T) {
	cfg := KWDBConfig{
		SQLPort:    26257,
		HTTPPort:   8080,
		DataDir:    "/data",
		LogDir:     "/logs",
		Insecure:   true,
		ListenAddr: "0.0.0.0:26257",
		HTTPAddr:   "0.0.0.0:8080",
	}

	// Test that fields are set correctly
	if cfg.SQLPort != 26257 {
		t.Errorf("Expected SQLPort 26257, got %d", cfg.SQLPort)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("Expected HTTPPort 8080, got %d", cfg.HTTPPort)
	}
	if cfg.DataDir != "/data" {
		t.Errorf("Expected DataDir '/data', got '%s'", cfg.DataDir)
	}
	if cfg.LogDir != "/logs" {
		t.Errorf("Expected LogDir '/logs', got '%s'", cfg.LogDir)
	}
	if cfg.Insecure != true {
		t.Errorf("Expected Insecure true, got %v", cfg.Insecure)
	}
	if cfg.ListenAddr != "0.0.0.0:26257" {
		t.Errorf("Expected ListenAddr '0.0.0.0:26257', got '%s'", cfg.ListenAddr)
	}
	if cfg.HTTPAddr != "0.0.0.0:8080" {
		t.Errorf("Expected HTTPAddr '0.0.0.0:8080', got '%s'", cfg.HTTPAddr)
	}
}
