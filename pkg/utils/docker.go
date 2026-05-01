package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// CheckDocker checks if Docker is installed and running
func CheckDocker() bool {
	cmd := exec.Command("docker", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// CheckDockerCompose checks if Docker Compose is available
func CheckDockerCompose() bool {
	// Try docker compose first (newer)
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return true
	}

	// Try docker-compose (older)
	cmd = exec.Command("docker-compose", "version")
	if err := cmd.Run(); err == nil {
		return true
	}

	return false
}

// GetDockerComposeCmd returns the docker compose command to use
func GetDockerComposeCmd() []string {
	// Try docker compose first
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return []string{"docker", "compose"}
	}

	// Fallback to docker-compose
	return []string{"docker-compose"}
}

// IsPortAvailable checks if a port is available
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// RunDocker runs a docker command
func RunDocker(args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetDockerImage checks if a docker image exists
func GetDockerImage(imageName string) bool {
	cmd := exec.Command("docker", "image", "ls", "-q", imageName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// PullDockerImage pulls a docker image
func PullDockerImage(imageName string) error {
	fmt.Printf("Pulling docker image: %s\n", imageName)
	cmd := exec.Command("docker", "pull", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

const (
	// AliyunRegistry is the Aliyun Docker registry
	AliyunRegistry = "registry.cn-hangzhou.aliyuncs.com/kwdb"
	// DefaultRegistry is the default Docker Hub registry
	DefaultRegistry = "kwdb"
)

// RewriteDockerComposeImage rewrites image names in docker-compose file to use Aliyun registry
func RewriteDockerComposeImage(composeFile, targetRegistry string) error {
	content, err := os.ReadFile(composeFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		// Only replace image: kwdb/xxx to targetRegistry/xxx
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "image:") {
			// Replace kwdb/xxx with targetRegistry/xxx
			lines[i] = strings.ReplaceAll(line, "image: "+DefaultRegistry+"/", "image: "+targetRegistry+"/")
			lines[i] = strings.ReplaceAll(lines[i], "image: docker.io/"+DefaultRegistry+"/", "image: "+targetRegistry+"/")
		}
	}
	newContent := strings.Join(lines, "\n")

	// Write back
	err = os.WriteFile(composeFile, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("  Rewrote image to use registry: %s\n", targetRegistry)
	return nil
}

// RestoreDockerComposeImage restores the original docker-compose file
func RestoreDockerComposeImage(composeFile, originalContent string) error {
	return os.WriteFile(composeFile, []byte(originalContent), 0644)
}

// ReadDockerComposeFile reads the docker-compose file content
func ReadDockerComposeFile(composeFile string) (string, error) {
	content, err := os.ReadFile(composeFile)
	return string(content), err
}
