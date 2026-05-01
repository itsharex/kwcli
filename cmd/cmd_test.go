package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmdSetup(t *testing.T) {
	// Test that root command is set up correctly
	if rootCmd == nil {
		t.Fatal("Expected rootCmd to not be nil")
	}

	if rootCmd.Use != "kwcli" {
		t.Errorf("Expected Use 'kwcli', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected Short to not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Expected Long to not be empty")
	}

	if rootCmd.Version != version {
		t.Errorf("Expected Version '%s', got '%s'", version, rootCmd.Version)
	}
}

func TestVersionValue(t *testing.T) {
	if version == "" {
		t.Error("Expected version to not be empty")
	}
}

func TestRootCmdHasGlobalFlags(t *testing.T) {
	// Check that global flags are registered
	flags := rootCmd.PersistentFlags()

	// Check config flag
	if !flags.Lookup("config").Changed {
		configFlag := flags.Lookup("config")
		if configFlag == nil {
			t.Error("Expected 'config' flag to exist")
		}
	}

	// Check home flag
	homeFlag := flags.Lookup("home")
	if homeFlag == nil {
		t.Error("Expected 'home' flag to exist")
	}
}

func TestGetHomeDir(t *testing.T) {
	// Test with custom home
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	// Note: homeDir is a package variable, need to call init or test differently
	// Just verify the function doesn't panic
	_ = GetHomeDir()
}

func TestGetComponentsDir(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	// Just verify the function doesn't panic
	_ = GetComponentsDir()
}

func TestEnsureDirs(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "kwcli-test-ensure-"+t.Name())
	defer os.RemoveAll(testHome)
	os.Setenv("KWCLI_HOME", testHome)
	defer os.Unsetenv("KWCLI_HOME")

	// Save original homeDir
	origHomeDir := homeDir
	homeDir = testHome
	defer func() { homeDir = origHomeDir }()

	// Call ensureDirs
	ensureDirs()

	// Check directories were created
	expectedDirs := []string{
		testHome,
		filepath.Join(testHome, "components"),
		filepath.Join(testHome, "data"),
		filepath.Join(testHome, "bin"),
	}

	for _, dir := range expectedDirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Errorf("Expected directory '%s' to exist, got error: %v", dir, err)
		}
		if !info.IsDir() {
			t.Errorf("Expected '%s' to be a directory", dir)
		}
	}
}

func TestPlaygroundCmdSetup(t *testing.T) {
	// Find playground command
	var playgroundCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "playground" {
			playgroundCmd = cmd
			break
		}
	}

	if playgroundCmd == nil {
		t.Fatal("Expected playground command to exist")
	}

	// Check subcommands
	expectedSubCommands := []string{"install", "start", "stop", "status", "restart", "upgrade", "logs"}
	for _, name := range expectedSubCommands {
		found := false
		for _, cmd := range playgroundCmd.Commands() {
			if cmd.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to exist", name)
		}
	}
}

func TestKWDBCmdSetup(t *testing.T) {
	// Find kwdb command
	var kwdbCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "kwdb" {
			kwdbCmd = cmd
			break
		}
	}

	if kwdbCmd == nil {
		t.Fatal("Expected kwdb command to exist")
	}

	// Check subcommands
	expectedSubCommands := []string{"install", "start", "stop", "status", "restart", "logs", "config"}
	for _, name := range expectedSubCommands {
		found := false
		for _, cmd := range kwdbCmd.Commands() {
			if cmd.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to exist", name)
		}
	}
}

func TestSQLCmdSetup(t *testing.T) {
	// Find sql command
	var sqlCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "sql" {
			sqlCmd = cmd
			break
		}
	}

	if sqlCmd == nil {
		t.Fatal("Expected sql command to exist")
	}

	// Check flags
	flags := sqlCmd.Flags()
	if flags.Lookup("host") == nil {
		t.Error("Expected 'host' flag to exist")
	}
	if flags.Lookup("user") == nil {
		t.Error("Expected 'user' flag to exist")
	}
	if flags.Lookup("database") == nil {
		t.Error("Expected 'database' flag to exist")
	}
	if flags.Lookup("execute") == nil {
		t.Error("Expected 'execute' flag to exist")
	}
	if flags.Lookup("insecure") == nil {
		t.Error("Expected 'insecure' flag to exist")
	}
}

func TestRootCmdLongDescription(t *testing.T) {
	longDesc := rootCmd.Long

	// Check that it contains key information
	if !strings.Contains(longDesc, "KWCLI") {
		t.Error("Expected Long to contain 'KWCLI'")
	}
	if !strings.Contains(longDesc, "Go Version") {
		t.Error("Expected Long to contain 'Go Version'")
	}
	if !strings.Contains(longDesc, "component") {
		t.Error("Expected Long to contain 'component'")
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute doesn't panic with no args
	// It should show help
	buf := bytes.NewBufferString("")
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{})

	err := Execute()
	// Execute will return an error because of missing subcommand
	// But we just want to make sure it doesn't panic
	_ = err
}

func TestHelpCommand(t *testing.T) {
	// Test that help command exists
	var helpCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "help" {
			helpCmd = cmd
			break
		}
	}

	if helpCmd == nil {
		t.Error("Expected help command to exist")
	}
}

func TestCompletionCommand(t *testing.T) {
	// Test that completion command exists
	var completionCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			completionCmd = cmd
			break
		}
	}

	if completionCmd == nil {
		t.Error("Expected completion command to exist")
	}
}
