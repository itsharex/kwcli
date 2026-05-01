package utils

import (
	"os"
	"testing"
)

func TestCheckDocker(t *testing.T) {
	// This test checks if docker is available
	// The result depends on the test environment
	result := CheckDocker()

	// Just verify it returns a boolean
	_ = result
}

func TestCheckDockerCompose(t *testing.T) {
	// This test checks if docker compose is available
	result := CheckDockerCompose()

	// Just verify it returns a boolean
	_ = result
}

func TestGetDockerComposeCmd(t *testing.T) {
	cmd := GetDockerComposeCmd()

	if len(cmd) == 0 {
		t.Error("Expected non-empty docker compose command")
	}

	// Should contain either "docker" or "docker-compose"
	if cmd[0] != "docker" && cmd[0] != "docker-compose" {
		t.Errorf("Expected 'docker' or 'docker-compose', got '%s'", cmd[0])
	}
}

func TestGetDockerImage(t *testing.T) {
	// Test with a known image that might exist
	// This is environment dependent
	result := GetDockerImage("alpine:latest")

	// Just verify it returns a boolean
	_ = result
}

func TestPullDockerImage(t *testing.T) {
	// This test actually pulls an image, which may take time
	// and require network access. We'll test with a small image.
	// Skip in CI or if docker is not available

	if os.Getenv("CI") != "" {
		t.Skip("Skipping in CI environment")
	}

	// Try to pull a small test image
	err := PullDockerImage("alpine:latest")
	if err != nil {
		t.Logf("Could not pull alpine image (may be expected): %v", err)
	}
}

func TestIsPortAvailable(t *testing.T) {
	// Test with a port that's unlikely to be in use
	result := IsPortAvailable(19999)
	if !result {
		t.Log("Port 19999 appears to be in use, but this might be expected")
	}

	// On Linux, port 0 is special and binds to an ephemeral port
	// So we skip this test as it's platform-dependent
	// Test with port > 65535 (invalid)
	result = IsPortAvailable(70000)
	if result {
		t.Error("Port 70000 should not be available")
	}
}

func TestRunDocker(t *testing.T) {
	// Test running docker --version
	err := RunDocker("--version")
	if err != nil {
		t.Logf("Docker may not be available: %v", err)
	}
}

func TestDownloadFile(t *testing.T) {
	// This test requires network access
	// We'll test with a small file from a reliable source
	if os.Getenv("CI") != "" {
		t.Skip("Skipping network test in CI")
	}

	dest := t.TempDir() + "/testfile"

	// Try to download a small file
	err := DownloadFile("https://httpbin.org/robots.txt", dest)
	if err != nil {
		t.Logf("Download failed (may be network issue): %v", err)
		return
	}

	// Check file exists
	info, err := os.Stat(dest)
	if err != nil {
		t.Errorf("Expected file to exist, got error: %v", err)
		return
	}

	if info.Size() == 0 {
		t.Error("Expected non-empty file")
	}
}

func TestExtractTarGz(t *testing.T) {
	// Test with invalid file
	err := ExtractTarGz("/nonexistent/file.tar.gz", t.TempDir())
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
