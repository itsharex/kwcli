package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shawn0915/kwcli/pkg/config"
	"github.com/spf13/cobra"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Connect to KWDB database",
	Long:  "Connect to KWDB database using kwbase CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		connectSQL(args)
	},
}

var (
	sqlHost      string
	sqlUser      string
	sqlDatabase  string
	sqlExecute   string
	insecureMode bool
)

func init() {
	rootCmd.AddCommand(sqlCmd)

	sqlCmd.Flags().StringVar(&sqlHost, "host", "127.0.0.1", "KWDB server host")
	sqlCmd.Flags().StringVarP(&sqlUser, "user", "u", "root", "KWDB user")
	sqlCmd.Flags().StringVarP(&sqlDatabase, "database", "d", "defaultdb", "Database name")
	sqlCmd.Flags().StringVarP(&sqlExecute, "execute", "e", "", "Execute SQL and exit")
	sqlCmd.Flags().BoolVar(&insecureMode, "insecure", true, "Use insecure mode (no TLS)")
}

func connectSQL(args []string) {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb")

	// Check Docker mode
	dockerMarker := installDir + "/.docker_mode"
	if _, err := os.Stat(dockerMarker); err == nil {
		connectSQLDocker(args)
		return
	}

	// Check Binary mode
	binaryMarker := installDir + "/.binary_mode"
	if _, err := os.Stat(binaryMarker); err == nil {
		connectSQLBinary(args)
		return
	}

	// Try to find kwbase anywhere in components directory
	kwbasePath := findKwbaseBinary(installDir)
	if kwbasePath != "" {
		connectSQLWithKwbase(kwbasePath, args)
		return
	}

	fmt.Println("KWDB is not installed. Run `kwcli kwdb install` first.")
	os.Exit(1)
}

func findKwbaseBinary(searchDir string) string {
	var kwbasePath string
	filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
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

func connectSQLWithKwbase(kwbasePath string, args []string) {
	// Build kwbase sql command
	cmdArgs := []string{kwbasePath, "sql"}

	// Add flags
	if insecureMode {
		cmdArgs = append(cmdArgs, "--insecure")
	}

	if sqlHost != "" {
		cmdArgs = append(cmdArgs, "--host="+sqlHost)
	}

	if sqlUser != "" && sqlUser != "root" {
		cmdArgs = append(cmdArgs, "-u", sqlUser)
	}

	if sqlDatabase != "" && sqlDatabase != "defaultdb" {
		cmdArgs = append(cmdArgs, "-d", sqlDatabase)
	}

	if sqlExecute != "" {
		cmdArgs = append(cmdArgs, "-e", sqlExecute)
	}

	// Run kwbase sql
	execCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Dir = filepath.Dir(kwbasePath)

	if err := execCmd.Run(); err != nil {
		fmt.Printf("Failed to execute SQL: %v\n", err)
		os.Exit(1)
	}
}

func connectSQLDocker(args []string) {
	// Check if container is running
	checkCmd := exec.Command("docker", "ps", "--filter", "name=kwdb", "--format", "{{.Names}}")
	output, err := checkCmd.Output()
	if err != nil || string(output) == "" {
		fmt.Println("KWDB container is not running. Run `kwcli kwdb start` first.")
		os.Exit(1)
	}

	// Build docker exec command for kwbase sql
	var dockerArgs []string
	if sqlExecute != "" {
		// Execute single SQL and exit
		dockerArgs = []string{"exec", "kwdb", "./kwbase", "sql", "--insecure", "--host=127.0.0.1"}
		if sqlUser != "" && sqlUser != "root" {
			dockerArgs = append(dockerArgs, "-u", sqlUser)
		}
		if sqlDatabase != "" && sqlDatabase != "defaultdb" {
			dockerArgs = append(dockerArgs, "-d", sqlDatabase)
		}
		dockerArgs = append(dockerArgs, "-e", sqlExecute)
	} else {
		// Interactive mode
		dockerArgs = []string{"exec", "-it", "kwdb", "./kwbase", "sql", "--insecure", "--host=127.0.0.1"}
		if sqlUser != "" && sqlUser != "root" {
			dockerArgs = append(dockerArgs, "-u", sqlUser)
		}
		if sqlDatabase != "" && sqlDatabase != "defaultdb" {
			dockerArgs = append(dockerArgs, "-d", sqlDatabase)
		}
	}

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to connect to KWDB: %v\n", err)
		os.Exit(1)
	}
}

func connectSQLBinary(args []string) {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb")

	// Find kwbase binary
	kwbasePath := findKwbaseBinary(installDir)
	if kwbasePath == "" {
		fmt.Println("kwbase binary not found. Please reinstall KWDB.")
		os.Exit(1)
	}

	connectSQLWithKwbase(kwbasePath, args)
}
