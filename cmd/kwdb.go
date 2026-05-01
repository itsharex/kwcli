package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shawn0915/kwcli/pkg/config"
	"github.com/shawn0915/kwcli/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	kwdbSQLPort     = 26257
	kwdbHTTPPort    = 8080
	kwdbDockerImage = "kwdb/kwdb"
)

// kwdbCmd represents the kwdb command
var kwdbCmd = &cobra.Command{
	Use:   "kwdb",
	Short: "Manage KWDB community edition (single-node)",
	Long:  "Manage KWDB community edition (single-node)",
}

// kwdbInstallCmd installs KWDB
var kwdbInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install KWDB community edition",
	Run: func(cmd *cobra.Command, args []string) {
		installKWDB()
	},
}

// kwdbStartCmd starts KWDB service
var kwdbStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start KWDB service",
	Run: func(cmd *cobra.Command, args []string) {
		startKWDB()
	},
}

// kwdbStopCmd stops KWDB service
var kwdbStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop KWDB service",
	Run: func(cmd *cobra.Command, args []string) {
		stopKWDB()
	},
}

// kwdbRestartCmd restarts KWDB service
var kwdbRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart KWDB service",
	Run: func(cmd *cobra.Command, args []string) {
		stopKWDB()
		startKWDB()
	},
}

// kwdbStatusCmd checks KWDB status
var kwdbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check KWDB status",
	Run: func(cmd *cobra.Command, args []string) {
		statusKWDB()
	},
}

// kwdbLogsCmd shows KWDB logs
var kwdbLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View KWDB logs",
	Run: func(cmd *cobra.Command, args []string) {
		logsKWDB()
	},
}

// kwdbConfigCmd represents the config subcommand
var kwdbConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage KWDB configuration",
	Long:  "Manage KWDB configuration",
}

// kwdbConfigShowCmd shows current configuration
var kwdbConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showKWDBConfig()
	},
}

// kwdbConfigEditCmd edits configuration
var kwdbConfigEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration interactively",
	Run: func(cmd *cobra.Command, args []string) {
		editKWDBConfig()
	},
}

// kwdbConfigSetCmd sets a configuration value
var kwdbConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setKWDBConfig(args[0], args[1])
	},
}

// kwdbConfigPathCmd shows config file path
var kwdbConfigPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show config file path",
	Run: func(cmd *cobra.Command, args []string) {
		showKWDBConfigPath()
	},
}

var nonInteractive bool

func init() {
	rootCmd.AddCommand(kwdbCmd)
	kwdbCmd.AddCommand(kwdbInstallCmd)
	kwdbCmd.AddCommand(kwdbStartCmd)
	kwdbCmd.AddCommand(kwdbStopCmd)
	kwdbCmd.AddCommand(kwdbRestartCmd)
	kwdbCmd.AddCommand(kwdbStatusCmd)
	kwdbCmd.AddCommand(kwdbLogsCmd)
	kwdbCmd.AddCommand(kwdbConfigCmd)

	kwdbConfigCmd.AddCommand(kwdbConfigShowCmd)
	kwdbConfigCmd.AddCommand(kwdbConfigEditCmd)
	kwdbConfigCmd.AddCommand(kwdbConfigSetCmd)
	kwdbConfigCmd.AddCommand(kwdbConfigPathCmd)

	kwdbInstallCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Use default configuration without prompting")
}

func installKWDB() {
	// Get or create configuration
	configPath := config.GetKWDBConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if !nonInteractive {
			fmt.Print("Would you like to configure KWDB? [Y/n]: ")
			var answer string
			fmt.Scanln(&answer)
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer == "" || answer == "y" || answer == "yes" {
				editKWDBConfig()
			} else {
				// Save default config
				cfg := config.DefaultKWDBConfig
				cfg.DataDir = filepath.Join(config.GetComponentsDir(), "kwdb", "data")
				cfg.LogDir = filepath.Join(config.GetComponentsDir(), "kwdb", "logs")
				config.SaveKWDBConfig(&cfg)
			}
		} else {
			cfg := config.DefaultKWDBConfig
			cfg.DataDir = filepath.Join(config.GetComponentsDir(), "kwdb", "data")
			cfg.LogDir = filepath.Join(config.GetComponentsDir(), "kwdb", "logs")
			config.SaveKWDBConfig(&cfg)
		}
	} else {
		fmt.Printf("Using existing configuration from %s\n", configPath)
	}

	// Try to install via binary
	if installKWDBBinary() {
		return
	}

	// Fallback to Docker
	fmt.Println("\nUnable to download the KWDB binary package.")
	fmt.Println("   This may be due to network issues or no package for your OS.")

	fmt.Print("Would you like to start KWDB using Docker instead? [Y/n]: ")
	var answer string
	fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer == "" || answer == "y" || answer == "yes" {
		installKWDBDocker()
	} else {
		fmt.Println("Installation cancelled.")
	}
}

func detectOSInfo() (string, string) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Map arch
	if arch == "x86_64" {
		arch = "x86_64"
	} else if arch == "arm64" {
		arch = "arm_64"
	}

	// Map OS
	osName := "linux"
	if os == "darwin" {
		osName = "darwin"
	} else if os == "windows" {
		osName = "windows"
	}

	return osName, arch
}

func installKWDBBinary() bool {
	osName, arch := detectOSInfo()
	fmt.Printf("Detected OS: %s, Arch: %s\n", osName, arch)

	// TODO: Fetch from Gitee API to get release info
	// For now, just return false to trigger Docker fallback
	fmt.Println("Binary installation from Gitee is not yet implemented.")
	return false
}

func installKWDBDocker() {
	if !utils.CheckDocker() {
		fmt.Println("Docker is not installed or not running.")
		os.Exit(1)
	}

	fmt.Println("Starting KWDB using Docker...")

	// Check ports
	if !utils.IsPortAvailable(kwdbSQLPort) {
		fmt.Printf("Port %d is already in use.\n", kwdbSQLPort)
		os.Exit(1)
	}
	if !utils.IsPortAvailable(kwdbHTTPPort) {
		fmt.Printf("Port %d is already in use.\n", kwdbHTTPPort)
		os.Exit(1)
	}

	// Pull image if not exists
	if !utils.GetDockerImage(kwdbDockerImage) {
		if err := utils.PullDockerImage(kwdbDockerImage); err != nil {
			fmt.Printf("Failed to pull image: %v\n", err)
			os.Exit(1)
		}
	}

	// Run container
	cmd := exec.Command("docker", "run", "-d", "--privileged", "--name", "kwdb",
		"-p", fmt.Sprintf("%d:%d", kwdbSQLPort, kwdbSQLPort),
		"-p", fmt.Sprintf("%d:%d", kwdbHTTPPort, kwdbHTTPPort),
		"-v", "/var/lib/kaiwudb:/kaiwudb/deploy/kaiwudb-container",
		"--ipc", "shareable",
		"-w", "/kaiwudb/bin",
		kwdbDockerImage,
		"./kwbase", "start-single-node",
		"--insecure",
		"--listen-addr=0.0.0.0:26257",
		"--http-addr=0.0.0.0:8080",
		"--store=/kaiwudb/deploy/kaiwudb-container",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to start KWDB container: %v\n", err)
		os.Exit(1)
	}

	// Write marker file
	markerPath := filepath.Join(config.GetComponentsDir(), "kwdb", ".docker_mode")
	os.MkdirAll(filepath.Dir(markerPath), 0755)
	os.WriteFile(markerPath, []byte("docker"), 0644)

	fmt.Println("\nKWDB (Docker) started successfully!")
	fmt.Printf("   SQL Port:  %d\n", kwdbSQLPort)
	fmt.Printf("   HTTP Port: %d\n", kwdbHTTPPort)
	fmt.Println("\nConnect using:")
	fmt.Println("   kwbase sql --insecure --host=127.0.0.1:26257")
}

func startKWDB() {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb")

	// Check Docker mode
	dockerMarker := installDir + "/.docker_mode"
	if _, err := os.Stat(dockerMarker); err == nil {
		fmt.Println("Starting KWDB Docker container...")
		cmd := exec.Command("docker", "start", "kwdb")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Check Binary mode
	binaryMarker := installDir + "/.binary_mode"
	if _, err := os.Stat(binaryMarker); err == nil {
		fmt.Println("Starting KWDB (binary mode)...")

		// Find kwbase binary
		kwbasePath := ""
		filepath.Walk(installDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (info.Name() == "kwbase" || info.Name() == "kwbase.exe") {
				kwbasePath = path
			}
			return nil
		})

		if kwbasePath == "" {
			fmt.Println("kwbase binary not found.")
			return
		}

		// Load config
		cfg, err := config.LoadKWDBConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			return
		}

		dataDir := cfg.DataDir
		if dataDir == "" {
			dataDir = filepath.Join(installDir, "data")
		}
		os.MkdirAll(dataDir, 0755)

		// Build start command
		cmdArgs := []string{kwbasePath, "start-single-node", "--store=" + dataDir}
		if cfg.Insecure {
			cmdArgs = append(cmdArgs, "--insecure")
		}
		if cfg.ListenAddr != "" {
			cmdArgs = append(cmdArgs, "--listen-addr="+cfg.ListenAddr)
		}
		if cfg.HTTPAddr != "" {
			cmdArgs = append(cmdArgs, "--http-addr="+cfg.HTTPAddr)
		}

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = filepath.Dir(kwbasePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	fmt.Println("KWDB is not installed. Run `kwcli kwdb install` first.")
}

func stopKWDB() {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb")

	// Check Docker mode
	dockerMarker := installDir + "/.docker_mode"
	if _, err := os.Stat(dockerMarker); err == nil {
		fmt.Println("Stopping KWDB Docker container...")
		cmd := exec.Command("docker", "stop", "kwdb")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Check Binary mode
	binaryMarker := installDir + "/.binary_mode"
	if _, err := os.Stat(binaryMarker); err == nil {
		fmt.Println("Stopping KWDB (binary mode)...")
		cmd := exec.Command("pkill", "-f", "kwbase")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	fmt.Println("KWDB is not running or not installed.")
}

func statusKWDB() {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb")

	// Check Docker mode
	dockerMarker := installDir + "/.docker_mode"
	if _, err := os.Stat(dockerMarker); err == nil {
		cmd := exec.Command("docker", "ps", "--all", "--filter", "name=^/kwdb$")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Check Binary mode
	binaryMarker := installDir + "/.binary_mode"
	if _, err := os.Stat(binaryMarker); err == nil {
		fmt.Println("Checking KWDB (binary mode) status...")
		cmd := exec.Command("pgrep", "-a", "kwbase")
		output, _ := cmd.Output()
		if len(output) > 0 {
			fmt.Println("KWDB is running:")
			fmt.Println(string(output))
		} else {
			fmt.Println("KWDB is not running.")
		}
		return
	}

	fmt.Println("KWDB is not installed. Run `kwcli kwdb install` first.")
}

func logsKWDB() {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb")

	// Check Docker mode
	dockerMarker := installDir + "/.docker_mode"
	if _, err := os.Stat(dockerMarker); err == nil {
		cmd := exec.Command("docker", "logs", "-f", "kwdb")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Check Binary mode
	binaryMarker := installDir + "/.binary_mode"
	if _, err := os.Stat(binaryMarker); err == nil {
		// Try to find logs in data directory
		cfg, err := config.LoadKWDBConfig()
		if err == nil && cfg.LogDir != "" {
			logFiles, _ := filepath.Glob(filepath.Join(cfg.LogDir, "*.log"))
			if len(logFiles) > 0 {
				fmt.Printf("Found logs in %s:\n", cfg.LogDir)
				for _, f := range logFiles {
					if len(logFiles) > 5 {
						break
					}
					fmt.Printf("  %s\n", f)
				}
				fmt.Println("\nUse `tail -f <log_file>` to view logs.")
				return
			}
		}

		// Show process info as fallback
		cmd := exec.Command("pgrep", "-a", "kwbase")
		output, _ := cmd.Output()
		if len(output) > 0 {
			fmt.Println("KWDB is running but no log files found in data directory.")
			fmt.Println("Process info:")
			fmt.Println(string(output))
		} else {
			fmt.Println("KWDB is not running.")
		}
		return
	}

	fmt.Println("KWDB is not installed. Run `kwcli kwdb install` first.")
}

func showKWDBConfig() {
	cfg, err := config.LoadKWDBConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	fmt.Println("Current KWDB Configuration:")
	fmt.Println("----------------------------------------")
	fmt.Printf("  sql_port: %d\n", cfg.SQLPort)
	fmt.Printf("  http_port: %d\n", cfg.HTTPPort)
	fmt.Printf("  data_dir: %s\n", cfg.DataDir)
	fmt.Printf("  log_dir: %s\n", cfg.LogDir)
	fmt.Printf("  insecure: %t\n", cfg.Insecure)
	fmt.Printf("  listen_addr: %s\n", cfg.ListenAddr)
	fmt.Printf("  http_addr: %s\n", cfg.HTTPAddr)
}

func editKWDBConfig() {
	cfg := config.DefaultKWDBConfig
	cfg.DataDir = filepath.Join(config.GetComponentsDir(), "kwdb", "data")
	cfg.LogDir = filepath.Join(config.GetComponentsDir(), "kwdb", "logs")

	// Interactive prompts
	fmt.Println("\n=== KWDB Configuration Wizard ===")
	fmt.Println("Press Enter to use the default value shown in brackets.")

	// SQL Port
	fmt.Printf("SQL Port [%d]: ", cfg.SQLPort)
	var portStr string
	fmt.Scanln(&portStr)
	if portStr != "" {
		fmt.Sscanf(portStr, "%d", &cfg.SQLPort)
	}

	// HTTP Port
	fmt.Printf("HTTP Port [%d]: ", cfg.HTTPPort)
	var httpPortStr string
	fmt.Scanln(&httpPortStr)
	if httpPortStr != "" {
		fmt.Sscanf(httpPortStr, "%d", &cfg.HTTPPort)
	}

	// Data Directory
	fmt.Printf("Data Directory [%s]: ", cfg.DataDir)
	var dataDir string
	fmt.Scanln(&dataDir)
	if dataDir != "" {
		cfg.DataDir = dataDir
	}

	// Log Directory
	fmt.Printf("Log Directory [%s]: ", cfg.LogDir)
	var logDir string
	fmt.Scanln(&logDir)
	if logDir != "" {
		cfg.LogDir = logDir
	}

	// Insecure mode
	fmt.Printf("Use insecure mode (disable authentication) [Y/n]: ")
	var insecureStr string
	fmt.Scanln(&insecureStr)
	cfg.Insecure = strings.ToLower(insecureStr) != "n"

	// Listen address
	fmt.Printf("SQL Listen Address [%s]: ", cfg.ListenAddr)
	var listenAddr string
	fmt.Scanln(&listenAddr)
	if listenAddr != "" {
		cfg.ListenAddr = listenAddr
	}

	// HTTP address
	fmt.Printf("HTTP Listen Address [%s]: ", cfg.HTTPAddr)
	var httpAddr string
	fmt.Scanln(&httpAddr)
	if httpAddr != "" {
		cfg.HTTPAddr = httpAddr
	}

	if err := config.SaveKWDBConfig(&cfg); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Println("Configuration saved. Restart KWDB for changes to take effect.")
}

func setKWDBConfig(key, value string) {
	cfg, err := config.LoadKWDBConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Set value based on key
	switch key {
	case "sql_port":
		fmt.Sscanf(value, "%d", &cfg.SQLPort)
	case "http_port":
		fmt.Sscanf(value, "%d", &cfg.HTTPPort)
	case "data_dir":
		cfg.DataDir = value
	case "log_dir":
		cfg.LogDir = value
	case "insecure":
		cfg.Insecure = value == "true"
	case "listen_addr":
		cfg.ListenAddr = value
	case "http_addr":
		cfg.HTTPAddr = value
	default:
		fmt.Printf("Unknown key: %s\n", key)
		return
	}

	if err := config.SaveKWDBConfig(cfg); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Printf("Set %s = %s\n", key, value)
	fmt.Println("Restart KWDB for changes to take effect.")
}

func showKWDBConfigPath() {
	configPath := config.GetKWDBConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Configuration file not found. Run `kwcli kwdb install` first.")
	} else {
		fmt.Println(configPath)
	}
}
