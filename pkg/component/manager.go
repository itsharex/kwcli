package component

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/KWDB/kwcli/pkg/config"
)

// Component represents a KWDB component
type Component struct {
	Name         string
	RepoURL      string
	RepoURLAlt   string // Alternative repository URL (e.g., AtomGit mirror)
	Description  string
	InstallDir   string
	Version      string // Current version (if installed)
}

// VersionMeta holds version information
type VersionMeta struct {
	Version   string
	Installed bool
}

// Source represents the repository source
type Source string

const (
	SourceAuto    Source = "auto"     // Auto-detect (try GitHub first, then fallback)
	SourceGitHub  Source = "github"   // GitHub
	SourceAtomGit Source = "atomgit"  // AtomGit (mirror)
)

// Registry holds all available components
var Registry = map[string]*Component{
	"playground": {
		Name:        "playground",
		RepoURL:     "https://github.com/KWDB/playground.git",
		RepoURLAlt:  "https://atomgit.com/kwdb/playground.git",
		Description: "KWDB Playground interactive learning platform",
	},
	"kwdb": {
		Name:        "kwdb",
		RepoURL:     "https://gitee.com/kwdb/kwdb",
		Description: "KWDB community edition (single-node)",
	},
}

// GetComponent returns a component by name
func GetComponent(name string) *Component {
	return Registry[name]
}

// ListComponents returns all available components
func ListComponents() []*Component {
	components := make([]*Component, 0, len(Registry))
	for _, c := range Registry {
		components = append(components, c)
	}
	return components
}

// IsInstalled checks if the component is installed (any version)
func (c *Component) IsInstalled() bool {
	versions := c.ListInstalledVersions()
	return len(versions) > 0
}

// IsVersionInstalled checks if a specific version is installed
func (c *Component) IsVersionInstalled(version string) bool {
	versionDir := c.GetVersionDir(version)
	info, err := os.Stat(versionDir)
	if err != nil {
		return false
	}
	return info.IsDir() && hasFiles(versionDir)
}

// EnsureInstallDir sets the install directory for the component
func (c *Component) EnsureInstallDir() {
	c.EnsureInstallDirWithVersion("")
}

// EnsureInstallDirWithVersion sets the install directory for the component with specific version
func (c *Component) EnsureInstallDirWithVersion(version string) {
	if version == "" {
		// Check current version
		currentVersion := c.GetCurrentVersion()
		if currentVersion != "" {
			version = currentVersion
		} else {
			// Use latest version if available
			latestVersion := c.GetLatestVersion()
			if latestVersion != "" {
				version = latestVersion
			}
		}
	}

	if version != "" {
		c.InstallDir = c.GetVersionDir(version)
	} else {
		c.InstallDir = filepath.Join(config.GetComponentsDir(), c.Name)
	}
	c.Version = version
}

// Install installs the component via git
func (c *Component) Install() error {
	return c.InstallWithVersionAndSource("", SourceAuto)
}

// InstallWithVersion installs the component with a specific version
func (c *Component) InstallWithVersion(version string) error {
	return c.InstallWithVersionAndSource(version, SourceAuto)
}

// InstallWithVersionAndSource installs the component with version and source
func (c *Component) InstallWithVersionAndSource(version string, source Source) error {
	// If no version specified and versions exist, use current version
	if version == "" {
		currentVersion := c.GetCurrentVersion()
		if currentVersion != "" {
			version = currentVersion
		} else {
			latestVersion := c.GetLatestVersion()
			if latestVersion != "" {
				version = latestVersion
			}
		}
	}

	// If still no version, use "latest" as default for new installations
	if version == "" {
		version = "latest"
	}

	// Set version-specific install directory
	c.EnsureInstallDirWithVersion(version)

	// Check if this specific version is already installed
	if c.IsVersionInstalled(version) {
		fmt.Printf("Component %s version %s is already installed at %s\n", c.Name, version, c.InstallDir)
		// Still set as current version
		c.SetCurrentVersion(version)
		return nil
	}

	// Create version directory
	parentDir := filepath.Dir(c.InstallDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed: %w", err)
	}

	fmt.Printf("Installing component %s...\n", c.Name)
	if version != "latest" {
		fmt.Printf("  Version: %s\n", version)
	}
	fmt.Printf("  Target: %s\n", c.InstallDir)

	// Try to clone with source selection
	err := c.tryClone(version, source)
	if err != nil && source == SourceAuto && c.RepoURLAlt != "" {
		// Auto mode: try fallback source
		fmt.Println("Primary source failed, trying alternative source...")
		err = c.tryClone(version, SourceAtomGit)
	}

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// If version was specified, save it for reference
	if version != "" {
		versionFile := filepath.Join(c.InstallDir, ".version")
		os.WriteFile(versionFile, []byte(version), 0644)
	}

	// Set current version to this version
	c.SetCurrentVersion(version)

	fmt.Printf("Component %s installed successfully!\n", c.Name)
	return nil
}

// tryClone attempts to clone the repository with the given source
func (c *Component) tryClone(version string, source Source) error {
	// Determine which repo URL to use
	var repoURL string
	var sourceName string

	switch source {
	case SourceGitHub:
		repoURL = c.RepoURL
		sourceName = "GitHub"
	case SourceAtomGit:
		if c.RepoURLAlt == "" {
			// Fallback to primary repo if no alternative is configured
			repoURL = c.RepoURL
			sourceName = c.RepoURL + " (fallback, no AtomGit mirror)"
		} else {
			repoURL = c.RepoURLAlt
			sourceName = "AtomGit"
		}
	default:
		// Auto mode: try primary repo first
		repoURL = c.RepoURL
		sourceName = "Auto"
	}

	fmt.Printf("  Repository: %s\n", repoURL)
	fmt.Printf("  Source: %s\n", sourceName)
	if version != "" {
		fmt.Printf("  Version: %s\n", version)
	}

	// Clone repository
	var cmd *exec.Cmd
	if version != "" && version != "latest" {
		// Clone specific tag/branch (full clone to ensure tag checkout works reliably)
		cmd = exec.Command("git", "clone", "-b", version, "--single-branch", repoURL, c.InstallDir)
	} else {
		// Shallow clone for latest/default branch
		cmd = exec.Command("git", "clone", "--depth=1", repoURL, c.InstallDir)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// If clone succeeded with a version, verify we're on the correct tag/branch
	if err == nil && version != "" && version != "latest" {
		// Check current branch/tag
		checkCmd := exec.Command("git", "describe", "--tags", "--exact-match", "--abbrev=0")
		checkCmd.Dir = c.InstallDir
		output, checkErr := checkCmd.Output()

		currentVersion := ""
		if checkErr == nil {
			currentVersion = strings.TrimSpace(string(output))
		}

		// If not on the expected version, try to checkout
		if currentVersion != version {
			fmt.Printf("  Verifying tag: %s\n", version)
			fetchCmd := exec.Command("git", "fetch", "origin", "tag", version, "--no-tags")
			fetchCmd.Dir = c.InstallDir
			fetchCmd.Run()

			checkoutCmd := exec.Command("git", "checkout", "-q", version)
			checkoutCmd.Dir = c.InstallDir
			if checkoutErr := checkoutCmd.Run(); checkoutErr != nil {
				return fmt.Errorf("version %s not found. Please check if the tag exists.", version)
			}
		}

		// Save version for reference
		versionFile := filepath.Join(c.InstallDir, ".version")
		os.WriteFile(versionFile, []byte(version), 0644)
	} else if err == nil && version != "" {
		// Save version for reference
		versionFile := filepath.Join(c.InstallDir, ".version")
		os.WriteFile(versionFile, []byte(version), 0644)
	}

	if err == nil {
		fmt.Printf("Component %s installed successfully!\n", c.Name)
	}
	return err
}

// Update updates the component to the latest version
func (c *Component) Update() error {
	return c.UpdateWithSource(SourceAuto)
}

// UpdateWithSource updates the component from a specific source
func (c *Component) UpdateWithSource(source Source) error {
	c.EnsureInstallDir()

	if !c.IsInstalled() {
		return c.InstallWithVersionAndSource("", source)
	}

	// Check if it's a git repository
	gitDir := filepath.Join(c.InstallDir, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		// Not a git repository, suggest reinstall
		fmt.Printf("Component %s was not installed via git, cannot update.\n", c.Name)
		fmt.Printf("Run 'kwcli %s uninstall' and then 'kwcli %s install' to reinstall.\n", c.Name, c.Name)
		return nil
	}

	// If source is specified and not auto, update remote URL if needed
	if source != SourceAuto {
		var expectedURL string
		switch source {
		case SourceGitHub:
			expectedURL = c.RepoURL
		case SourceAtomGit:
			expectedURL = c.RepoURLAlt
			if expectedURL == "" {
				expectedURL = c.RepoURL
			}
		}

		if expectedURL != "" {
			// Get current remote URL
			cmd := exec.Command("git", "remote", "get-url", "origin")
			cmd.Dir = c.InstallDir
			currentURLBytes, err := cmd.Output()
			currentURL := strings.TrimSpace(string(currentURLBytes))

			if err != nil || currentURL != expectedURL {
				fmt.Printf("Switching remote to %s...\n", expectedURL)
				setURLCmd := exec.Command("git", "remote", "set-url", "origin", expectedURL)
				setURLCmd.Dir = c.InstallDir
				if err := setURLCmd.Run(); err != nil {
					return fmt.Errorf("failed to update remote URL: %w", err)
				}
			}
		}
	}

	fmt.Printf("Updating component %s...\n", c.Name)

	// Git pull
	cmd := exec.Command("git", "pull")
	cmd.Dir = c.InstallDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update component: %w", err)
	}

	fmt.Printf("Component %s updated successfully!\n", c.Name)
	return nil
}

// Uninstall removes the component
func (c *Component) Uninstall() error {
	c.EnsureInstallDir()

	if !c.IsInstalled() {
		fmt.Printf("Component %s is not installed.\n", c.Name)
		return nil
	}

	fmt.Printf("Uninstalling component %s...\n", c.Name)

	if err := os.RemoveAll(c.InstallDir); err != nil {
		return fmt.Errorf("failed to remove component: %w", err)
	}

	fmt.Printf("Component %s uninstalled successfully!\n", c.Name)
	return nil
}

// GetInstallStatus returns the installation status of a component
func GetInstallStatus(name string) string {
	c := GetComponent(name)
	if c == nil {
		return "unknown"
	}
	if c.IsInstalled() {
		return "installed"
	}
	return "not installed"
}

// GetVersionDir returns the version-specific install directory
func (c *Component) GetVersionDir(version string) string {
	return filepath.Join(config.GetComponentsDir(), c.Name, "versions", version)
}

// GetVersionsDir returns the versions directory
func (c *Component) GetVersionsDir() string {
	return filepath.Join(config.GetComponentsDir(), c.Name, "versions")
}

// GetCurrentVersion returns the currently active version
func (c *Component) GetCurrentVersion() string {
	versionsDir := c.GetVersionsDir()
	currentFile := filepath.Join(versionsDir, ".current")
	data, err := os.ReadFile(currentFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// SetCurrentVersion sets the currently active version
func (c *Component) SetCurrentVersion(version string) error {
	versionsDir := c.GetVersionsDir()
	os.MkdirAll(versionsDir, 0755)
	currentFile := filepath.Join(versionsDir, ".current")
	return os.WriteFile(currentFile, []byte(version), 0644)
}

// ListInstalledVersions returns all installed versions
func (c *Component) ListInstalledVersions() []string {
	versionsDir := c.GetVersionsDir()
	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return []string{}
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			versionDir := filepath.Join(versionsDir, entry.Name())
			if hasFiles(versionDir) {
				versions = append(versions, entry.Name())
			}
		}
	}
	return versions
}

// GetLatestVersion returns the latest installed version
func (c *Component) GetLatestVersion() string {
	versions := c.ListInstalledVersions()
	if len(versions) == 0 {
		return ""
	}
	// Return the last one (assumes versions are sorted)
	return versions[len(versions)-1]
}

// hasFiles checks if a directory has any files
func hasFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// GetComposeFile returns the path to docker-compose.yml
func (c *Component) GetComposeFile(subPath string) string {
	return filepath.Join(c.InstallDir, subPath, "docker-compose.yml")
}

// RunDockerCompose runs docker compose command
func (c *Component) RunDockerCompose(args ...string) error {
	cmd := exec.Command("docker", append([]string{"compose", "-f", c.GetComposeFile("docker/playground")}, args...)...)
	cmd.Dir = c.InstallDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetBinaryPath returns the path to the kwbase binary
func (c *Component) GetBinaryPath() string {
	// Try to find kwbase in the installed directory
	var kwbasePath string
	
	filepath.Walk(c.InstallDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && (info.Name() == "kwbase" || info.Name() == "kwbase.exe") {
			kwbasePath = path
		}
		return nil
	})
	
	return kwbasePath
}

// Get kwdb release info from Gitee
func GetGiteeLatestRelease() (string, error) {
	// This would require HTTP client - simplified for now
	return "", nil
}

// FindMatchingAsset finds a matching binary asset for the current OS/arch
func FindMatchingAsset(releaseAssets []string, osName, arch string) string {
	for _, asset := range releaseAssets {
		lowerAsset := strings.ToLower(asset)
		lowerOS := strings.ToLower(osName)
		lowerArch := strings.ToLower(arch)
		
		if strings.Contains(lowerAsset, lowerOS) && strings.Contains(lowerAsset, lowerArch) {
			return asset
		}
	}
	return ""
}