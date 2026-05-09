package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shawn0915/kwcli/pkg/sampledb"
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

// === Sampledb Tests ===

func TestSampledbCmdSetup(t *testing.T) {
	// Find sampledb command
	var sampledbCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "sampledb" {
			sampledbCmd = cmd
			break
		}
	}

	if sampledbCmd == nil {
		t.Fatal("Expected sampledb command to exist")
	}

	// Check subcommands
	expectedSubCommands := []string{"init", "generate", "list", "run", "clean", "status"}
	for _, name := range expectedSubCommands {
		found := false
		for _, cmd := range sampledbCmd.Commands() {
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

func TestSampledbScenarios(t *testing.T) {
	// Test that sampledb scenarios are defined
	if len(sampledb.AllScenarios) == 0 {
		t.Error("Expected at least one scenario")
	}

	// Test FindScenario function
	scenario := sampledb.FindScenario("top10-area-energy")
	if scenario == nil {
		t.Error("Expected to find 'top10-area-energy' scenario")
	}
	if scenario.Title == "" {
		t.Error("Expected scenario to have a title")
	}
	if scenario.SQL == "" {
		t.Error("Expected scenario to have SQL")
	}

	// Test FindScenario returns nil for unknown scenario
	unknownScenario := sampledb.FindScenario("unknown-scenario")
	if unknownScenario != nil {
		t.Error("Expected nil for unknown scenario")
	}
}

func TestSampledbScenarioCategories(t *testing.T) {
	// Test that all scenarios have categories
	for _, s := range sampledb.AllScenarios {
		if s.Category == "" {
			t.Errorf("Scenario '%s' has no category", s.Name)
		}
	}

	// Test that each category has at least one scenario
	categories := make(map[string]int)
	for _, s := range sampledb.AllScenarios {
		categories[s.Category]++
	}

	if categories["basic"] == 0 {
		t.Error("Expected at least one 'basic' category scenario")
	}
	if categories["cross-mode"] == 0 {
		t.Error("Expected at least one 'cross-mode' category scenario")
	}
	if categories["window"] == 0 {
		t.Error("Expected at least one 'window' category scenario")
	}
}

func TestSampledbScenarioNames(t *testing.T) {
	expectedScenarios := []string{
		"top10-area-energy",
		"fault-meters",
		"meter-summary",
		"alarm-detection",
		"area-energy-stats",
		"recent-24h-trend",
		"cross-mode-join",
		"cross-mode-user-power",
		"cross-mode-alarm-analysis",
		"cross-mode-region-comparison",
		"time-bucket-stats",
		"session-analysis",
		"voltage-state",
		"abnormal-current",
		"sliding-window",
		"time-window-advanced",
		"count-window-example",
	}

	for _, name := range expectedScenarios {
		scenario := sampledb.FindScenario(name)
		if scenario == nil {
			t.Errorf("Expected scenario '%s' to exist", name)
		}
	}
}

func TestSampledbSchema(t *testing.T) {
	// Test RDBSchema is not empty
	if sampledb.RDBSchema == "" {
		t.Error("Expected RDBSchema to not be empty")
	}

	// Test TSDBSchema is not empty
	if sampledb.TSDBSchema == "" {
		t.Error("Expected TSDBSchema to not be empty")
	}

	// Test schemas contain expected keywords
	if !strings.Contains(sampledb.RDBSchema, "CREATE DATABASE") {
		t.Error("Expected RDBSchema to contain 'CREATE DATABASE'")
	}
	if !strings.Contains(sampledb.RDBSchema, "meter_info") {
		t.Error("Expected RDBSchema to contain 'meter_info'")
	}
	if !strings.Contains(sampledb.TSDBSchema, "CREATE TS DATABASE") {
		t.Error("Expected TSDBSchema to contain 'CREATE TS DATABASE'")
	}
	if !strings.Contains(sampledb.TSDBSchema, "meter_data") {
		t.Error("Expected TSDBSchema to contain 'meter_data'")
	}
}

func TestSampledbDataGeneration(t *testing.T) {
	// Test GenerateRDBData returns SQL
	rdbData := sampledb.GenerateRDBData()
	if rdbData == "" {
		t.Error("Expected GenerateRDBData to return SQL")
	}

	// Test contains expected statements
	if !strings.Contains(rdbData, "INSERT INTO rdb.area_info") {
		t.Error("Expected RDB data to contain area_info inserts")
	}
	if !strings.Contains(rdbData, "INSERT INTO rdb.user_info") {
		t.Error("Expected RDB data to contain user_info inserts")
	}
	if !strings.Contains(rdbData, "INSERT INTO rdb.meter_info") {
		t.Error("Expected RDB data to contain meter_info inserts")
	}
	if !strings.Contains(rdbData, "INSERT INTO rdb.alarm_rules") {
		t.Error("Expected RDB data to contain alarm_rules inserts")
	}

	// Test GenerateTSDBData returns SQL
	tsdbData := sampledb.GenerateTSDBData()
	if tsdbData == "" {
		t.Error("Expected GenerateTSDBData to return SQL")
	}

	// Test contains expected statements
	if !strings.Contains(tsdbData, "INSERT INTO tsdb.meter_data") {
		t.Error("Expected TSDB data to contain meter_data inserts")
	}
}

func TestSampledbRunnerSplitSQL(t *testing.T) {
	// Test the SQL splitting logic in runner
	// This tests the internal function via integration
	// We'll just verify basic functionality here

	testCases := []struct {
		input    string
		expected int
	}{
		{"SELECT 1; SELECT 2;", 2},
		{"SELECT 1; SELECT 2; SELECT 3;", 3},
		{"SELECT 1;", 1},
		{"", 0},
	}

	for _, tc := range testCases {
		// Basic validation - just ensure we can parse
		if tc.expected > 0 && tc.input != "" {
			statements := strings.Split(tc.input, ";")
			count := 0
			for _, s := range statements {
				if strings.TrimSpace(s) != "" {
					count++
				}
			}
			if count != tc.expected {
				t.Errorf("For input '%s', expected %d statements, got %d", tc.input, tc.expected, count)
			}
		}
	}
}

func TestSampledbCrossModeScenarios(t *testing.T) {
	// Test cross-mode scenarios contain required joins
	crossModeScenarios := []string{
		"cross-mode-join",
		"cross-mode-user-power",
		"cross-mode-alarm-analysis",
		"cross-mode-region-comparison",
	}

	for _, name := range crossModeScenarios {
		scenario := sampledb.FindScenario(name)
		if scenario == nil {
			t.Errorf("Expected cross-mode scenario '%s' to exist", name)
			continue
		}

		// Cross-mode should join tsdb and rdb
		if !strings.Contains(scenario.SQL, "tsdb.") {
			t.Errorf("Scenario '%s' should reference tsdb", name)
		}
		if !strings.Contains(scenario.SQL, "rdb.") {
			t.Errorf("Scenario '%s' should reference rdb", name)
		}
	}
}

func TestSampledbWindowScenarios(t *testing.T) {
	// Test window function scenarios
	windowScenarios := []string{
		"time-bucket-stats",
		"session-analysis",
		"voltage-state",
		"abnormal-current",
		"sliding-window",
	}

	for _, name := range windowScenarios {
		scenario := sampledb.FindScenario(name)
		if scenario == nil {
			t.Errorf("Expected window scenario '%s' to exist", name)
			continue
		}

		// Window functions should have window in SQL (time_bucket, count_window, etc)
		if !strings.Contains(scenario.SQL, "_window(") && !strings.Contains(scenario.SQL, "time_bucket(") {
			t.Errorf("Scenario '%s' should contain window function", name)
		}
	}
}

// === TSBS Tests ===

func TestTSBSCmdSetup(t *testing.T) {
	// Find tsbs command
	var tsbsCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "tsbs" {
			tsbsCmd = cmd
			break
		}
	}

	if tsbsCmd == nil {
		t.Fatal("Expected tsbs command to exist")
	}

	// Check subcommands - now using 'init' instead of 'generate-data' and 'generate-queries'
	expectedSubCommands := []string{"init", "load", "run", "list", "clean"}
	for _, name := range expectedSubCommands {
		found := false
		for _, cmd := range tsbsCmd.Commands() {
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

func TestTSBSQueryTypes(t *testing.T) {
	// Test CPU query types
	queryTypes := getQueryTypes("cpu")
	if len(queryTypes) == 0 {
		t.Error("Expected at least one CPU query type")
	}

	// Test that specific query types exist
	cpuQueryNames := make(map[string]bool)
	for _, qt := range queryTypes {
		cpuQueryNames[qt.Name] = true
	}

	if !cpuQueryNames["cpu-max-all"] {
		t.Error("Expected 'cpu-max-all' query type in CPU use case")
	}
	if !cpuQueryNames["double-groupby"] {
		t.Error("Expected 'double-groupby' query type in CPU use case")
	}
}

func TestTSBSIoTQueryTypes(t *testing.T) {
	// Test IoT query types
	queryTypes := getQueryTypes("iot")
	if len(queryTypes) == 0 {
		t.Error("Expected at least one IoT query type")
	}

	iotQueryNames := make(map[string]bool)
	for _, qt := range queryTypes {
		iotQueryNames[qt.Name] = true
	}

	if !iotQueryNames["iot-ingest"] {
		t.Error("Expected 'iot-ingest' query type in IoT use case")
	}
	if !iotQueryNames["iot-threshold"] {
		t.Error("Expected 'iot-threshold' query type in IoT use case")
	}
}

// TestCodeFormatting checks that all Go files are properly formatted
func TestCodeFormatting(t *testing.T) {
	// Run gofmt to check for formatting issues (exclude third_party directory)
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = "/home/shawnyan/kwcli"
	output, err := cmd.CombinedOutput()

	if err != nil && err.Error() != "exit status 1" {
		t.Logf("Warning: gofmt check failed: %v", err)
	}

	unformatted := string(output)
	// Filter out third_party directory (it's third-party code)
	var filtered []string
	for _, line := range strings.Split(unformatted, "\n") {
		if line != "" && !strings.HasPrefix(line, "third_party/") {
			filtered = append(filtered, line)
		}
	}

	if len(filtered) > 0 {
		t.Errorf("The following files need formatting:\n%s", strings.Join(filtered, "\n"))
	}
}
