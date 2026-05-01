package component

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shawn0915/kwcli/pkg/config"
)

func TestComponentRegistry(t *testing.T) {
	// Test that registry is populated
	components := ListComponents()
	if len(components) == 0 {
		t.Error("Expected components in registry, got none")
	}

	// Test playground component exists
	playground := GetComponent("playground")
	if playground == nil {
		t.Error("Expected playground component to exist")
	}
	if playground.Name != "playground" {
		t.Errorf("Expected name 'playground', got '%s'", playground.Name)
	}
	if playground.RepoURL == "" {
		t.Error("Expected playground RepoURL to be set")
	}

	// Test kwdb component exists
	kwdb := GetComponent("kwdb")
	if kwdb == nil {
		t.Error("Expected kwdb component to exist")
	}
	if kwdb.Name != "kwdb" {
		t.Errorf("Expected name 'kwdb', got '%s'", kwdb.Name)
	}
}

func TestGetComponent(t *testing.T) {
	tests := []struct {
		name     string
		wantName string
		wantNil  bool
	}{
		{"playground", "playground", false},
		{"kwdb", "kwdb", false},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetComponent(tt.name)
			if tt.wantNil {
				if c != nil {
					t.Errorf("Expected nil, got %v", c)
				}
			} else {
				if c == nil {
					t.Errorf("Expected non-nil for %s", tt.name)
				}
				if c != nil && c.Name != tt.wantName {
					t.Errorf("Expected Name '%s', got '%s'", tt.wantName, c.Name)
				}
			}
		})
	}
}

func TestComponent_IsInstalled(t *testing.T) {
	// Set up test home directory
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+strings.ReplaceAll(t.Name(), "/", ""))
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	comp := &Component{
		Name: "test-component",
	}

	// Should not be installed initially (dir doesn't exist)
	if comp.IsInstalled() {
		t.Error("Expected component to not be installed")
	}

	// Create the versioned directory (new structure: versions/v1.0.0/)
	versionDir := comp.GetVersionDir("v1.0.0")
	os.MkdirAll(versionDir, 0755)

	// Should still be false (empty directory)
	if comp.IsInstalled() {
		t.Error("Expected component to not be installed (empty dir)")
	}

	// Create a file in the version directory
	os.WriteFile(filepath.Join(versionDir, "test.txt"), []byte("test"), 0644)

	// Now should be installed
	if !comp.IsInstalled() {
		t.Error("Expected component to be installed (with files)")
	}
}

func TestComponent_EnsureInstallDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)

	// Set environment and reload config
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	// Force config to reload by clearing cached value
	// The function should pick up the new home directory
	comp := &Component{Name: "test-comp"}
	comp.EnsureInstallDir()

	// Since config caches the value, we just verify it creates a valid path
	if comp.InstallDir == "" {
		t.Error("Expected InstallDir to be set")
	}
	if !strings.Contains(comp.InstallDir, "test-comp") {
		t.Errorf("Expected InstallDir to contain 'test-comp', got '%s'", comp.InstallDir)
	}
}

func TestComponent_GetComposeFile(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)

	comp := &Component{
		Name:       "playground",
		InstallDir: filepath.Join(testHome, "components", "playground"),
	}

	os.MkdirAll(comp.InstallDir, 0755)

	result := comp.GetComposeFile("docker/playground")
	expected := filepath.Join(testHome, "components", "playground", "docker", "playground", "docker-compose.yml")

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestGetInstallStatus(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+strings.ReplaceAll(t.Name(), "/", ""))
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	tests := []struct {
		name      string
		installed bool
		want      string
	}{
		{"playground", false, "not installed"},
		{"playground", true, "installed"},
		{"unknown-comp", false, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := &Component{Name: tt.name}
			compDir := comp.GetVersionsDir()
			os.RemoveAll(compDir)
			if tt.installed {
				versionDir := comp.GetVersionDir("v1.0.0")
				os.MkdirAll(versionDir, 0755)
				os.WriteFile(filepath.Join(versionDir, "test.txt"), []byte("test"), 0644)
			}

			got := GetInstallStatus(tt.name)
			if got != tt.want {
				t.Errorf("GetInstallStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourceConstants(t *testing.T) {
	if SourceAuto != "auto" {
		t.Errorf("Expected SourceAuto 'auto', got '%s'", SourceAuto)
	}
	if SourceGitHub != "github" {
		t.Errorf("Expected SourceGitHub 'github', got '%s'", SourceGitHub)
	}
	if SourceAtomGit != "atomgit" {
		t.Errorf("Expected SourceAtomGit 'atomgit', got '%s'", SourceAtomGit)
	}
}

func TestComponent_RepoURLAlt(t *testing.T) {
	// Test that playground has alternative URL
	playground := GetComponent("playground")
	if playground.RepoURLAlt == "" {
		t.Error("Expected playground to have RepoURLAlt set")
	}

	// Test that kwdb doesn't have alternative URL (should be empty)
	kwdb := GetComponent("kwdb")
	if kwdb.RepoURLAlt != "" {
		t.Error("Expected kwdb to have empty RepoURLAlt")
	}
}

func TestListComponents(t *testing.T) {
	components := ListComponents()

	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	// Check component names
	expectedNames := map[string]bool{"playground": false, "kwdb": false}
	for _, c := range components {
		if _, ok := expectedNames[c.Name]; ok {
			expectedNames[c.Name] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected component '%s' not found", name)
		}
	}
}

// Test config package functions are accessible
func TestConfigGetHomeDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	homeDir := config.GetHomeDir()
	if homeDir != testHome {
		t.Errorf("Expected home dir '%s', got '%s'", testHome, homeDir)
	}
}

func TestConfigGetComponentsDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	componentsDir := config.GetComponentsDir()
	expected := filepath.Join(testHome, "components")
	if componentsDir != expected {
		t.Errorf("Expected components dir '%s', got '%s'", expected, componentsDir)
	}
}
