package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shawn0915/kwcli/pkg/config"
	"github.com/spf13/cobra"
)

// tsbsCmd represents the tsbs command
var tsbsCmd = &cobra.Command{
	Use:   "tsbs",
	Short: "KWDB Time-Series Benchmark Suite (TSBS)",
	Long: `KWDB Time-Series Benchmark Suite (TSBS)
A high-performance benchmarking tool for KWDB time-series database.

Commands:
  init            Initialize benchmark: generate data and queries
  load            Load data into KWDB
  run             Run query benchmarks
  list            List available query types

Examples:
  kwcli tsbs init --help
  kwcli tsbs load --help
  kwcli tsbs run --help`,
}

// Detect TSBS installation and path
func getTSBSPath() (string, error) {
	installDir := filepath.Join(config.GetComponentsDir(), "kwdb-tsbs")

	// Check Docker mode
	dockerMarker := installDir + "/.docker_mode"
	if _, err := os.Stat(dockerMarker); err == nil {
		return "docker", nil
	}

	// Check Binary mode
	binaryMarker := installDir + "/.binary_mode"
	if _, err := os.Stat(binaryMarker); err == nil {
		// Find tsbs binary
		tsbsPath := findTSBS(installDir)
		if tsbsPath != "" {
			return tsbsPath, nil
		}
	}

	// Try to find tsbs anywhere
	tsbsPath := findTSBS(installDir)
	if tsbsPath != "" {
		return tsbsPath, nil
	}

	return "", fmt.Errorf("kwdb-tsbs is not installed. Please run 'kwcli tsbs install' or clone from https://github.com/KWDB/kwdb-tsbs")
}

func findTSBS(searchDir string) string {
	var tsbsPath string
	filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && (info.Name() == "tsbs_generate_data" || info.Name() == "tsbs_generate_data.exe") {
			tsbsPath = filepath.Dir(path)
		}
		return nil
	})
	return tsbsPath
}

// TSBSInitCmd generates both benchmark data and queries in one command
var TSBSInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize benchmark: generate data and queries",
	Long: `Generate both benchmark test data and query files in one command.
This combines generate-data and generate-queries into a single operation.

Examples:
  # Initialize with default settings (CPU, 10 devices, 1 day)
  kwcli tsbs init

  # Initialize with custom settings
  kwcli tsbs init --use-case=cpu --scale=100 --timestamp-start="2024-01-01T00:00:00Z" --timestamp-end="2024-01-02T00:00:00Z" --queries=5000

  # Initialize for IoT use case
  kwcli tsbs init --use-case=iot --scale=50 --days=7 --queries=10000`,
	Run: func(cmd *cobra.Command, args []string) {
		tsbsPath, err := getTSBSPath()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if tsbsPath == "docker" {
			fmt.Println("Docker mode: use 'kwcli tsbs docker-init' instead")
			os.Exit(1)
		}

		// Get flag values
		useCase, _ := cmd.Flags().GetString("use-case")
		scale, _ := cmd.Flags().GetInt("scale")
		tsStart, _ := cmd.Flags().GetString("timestamp-start")
		tsEnd, _ := cmd.Flags().GetString("timestamp-end")
		interval, _ := cmd.Flags().GetString("sampling-interval")
		queries, _ := cmd.Flags().GetInt("queries")
		queryType, _ := cmd.Flags().GetString("query-type")
		dataFile, _ := cmd.Flags().GetString("data-file")
		queryFile, _ := cmd.Flags().GetString("query-file")
		outputDir, _ := cmd.Flags().GetString("output-dir")

		// Set default output paths
		if dataFile == "" {
			if outputDir != "" {
				dataFile = filepath.Join(outputDir, "tsbs_data")
				queryFile = filepath.Join(outputDir, "tsbs_queries")
			} else {
				dataFile = filepath.Join(os.TempDir(), "tsbs_data")
				queryFile = filepath.Join(os.TempDir(), "tsbs_queries")
			}
		}

		// ============================================
		// Step 1: Generate benchmark data
		// ============================================
		fmt.Println("=== Step 1/2: Generating benchmark data ===")

		dataBinaryPath := filepath.Join(tsbsPath, "tsbs_generate_data")
		if _, err := os.Stat(dataBinaryPath); os.IsNotExist(err) {
			dataBinaryPath = dataBinaryPath + ".exe"
		}

		dataCmdArgs := []string{dataBinaryPath}
		dataCmdArgs = append(dataCmdArgs, "--use-case", useCase)
		dataCmdArgs = append(dataCmdArgs, "--scale", fmt.Sprintf("%d", scale))
		dataCmdArgs = append(dataCmdArgs, "--timestamp-start", tsStart)
		dataCmdArgs = append(dataCmdArgs, "--timestamp-end", tsEnd)
		dataCmdArgs = append(dataCmdArgs, "--sampling-interval", interval)
		dataCmdArgs = append(dataCmdArgs, "--file", dataFile)

		fmt.Printf("Executing: %s\n", strings.Join(dataCmdArgs, " "))

		execDataCmd := exec.Command(dataCmdArgs[0], dataCmdArgs[1:]...)
		execDataCmd.Stdout = os.Stdout
		execDataCmd.Stderr = os.Stderr
		execDataCmd.Dir = tsbsPath
		if err := execDataCmd.Run(); err != nil {
			fmt.Printf("Error generating data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ Data generation completed!")
		fmt.Printf("  Data file: %s\n", dataFile)

		// ============================================
		// Step 2: Generate query files
		// ============================================
		fmt.Println("\n=== Step 2/2: Generating query files ===")

		queryBinaryPath := filepath.Join(tsbsPath, "tsbs_generate_queries")
		if _, err := os.Stat(queryBinaryPath); os.IsNotExist(err) {
			queryBinaryPath = queryBinaryPath + ".exe"
		}

		queryCmdArgs := []string{queryBinaryPath}
		queryCmdArgs = append(queryCmdArgs, "--use-case", useCase)
		queryCmdArgs = append(queryCmdArgs, "--scale", fmt.Sprintf("%d", scale))
		queryCmdArgs = append(queryCmdArgs, "--timestamp-start", tsStart)
		queryCmdArgs = append(queryCmdArgs, "--timestamp-end", tsEnd)
		queryCmdArgs = append(queryCmdArgs, "--queries", fmt.Sprintf("%d", queries))
		if queryType != "" {
			queryCmdArgs = append(queryCmdArgs, "--query-type", queryType)
		}
		queryCmdArgs = append(queryCmdArgs, "--file", queryFile)

		fmt.Printf("Executing: %s\n", strings.Join(queryCmdArgs, " "))

		execQueryCmd := exec.Command(queryCmdArgs[0], queryCmdArgs[1:]...)
		execQueryCmd.Stdout = os.Stdout
		execQueryCmd.Stderr = os.Stderr
		execQueryCmd.Dir = tsbsPath
		if err := execQueryCmd.Run(); err != nil {
			fmt.Printf("Error generating queries: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ Query generation completed!")
		fmt.Printf("  Query file: %s\n", queryFile)

		// ============================================
		// Summary
		// ============================================
		fmt.Println("\n=== Initialization Complete ===")
		fmt.Printf("Use Case:     %s\n", useCase)
		fmt.Printf("Scale:        %d devices\n", scale)
		fmt.Printf("Time Range:   %s to %s\n", tsStart, tsEnd)
		fmt.Printf("Queries:      %d\n", queries)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. Load data:      kwcli tsbs load --file=%s\n", dataFile)
		fmt.Printf("  2. Run queries:   kwcli tsbs run --file=%s\n", queryFile)
	},
}

// TSBSLoadCmd loads data into KWDB
var TSBSLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load data into KWDB",
	Long: `Load generated benchmark data into KWDB.

Examples:
  # Load data from file
  kwcli tsbs load --file=/tmp/tsbs_data

  # Load with batch size
  kwcli tsbs load --file=/tmp/tsbs_data --batch-size=10000`,
	Run: func(cmd *cobra.Command, args []string) {
		tsbsPath, err := getTSBSPath()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if tsbsPath == "docker" {
			fmt.Println("Docker mode: use 'kwcli tsbs docker-load' instead")
			os.Exit(1)
		}

		binaryPath := filepath.Join(tsbsPath, "tsbs_load_kwdb")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			binaryPath = binaryPath + ".exe"
		}

		// Build command
		cmdArgs := []string{binaryPath}

		// Host
		if host, _ := cmd.Flags().GetString("host"); host != "" {
			cmdArgs = append(cmdArgs, "--host", host)
		}

		// Port
		if port, _ := cmd.Flags().GetInt("port"); port > 0 {
			cmdArgs = append(cmdArgs, "--port", fmt.Sprintf("%d", port))
		}

		// User
		if user, _ := cmd.Flags().GetString("user"); user != "" {
			cmdArgs = append(cmdArgs, "--user", user)
		}

		// Password
		if password, _ := cmd.Flags().GetString("password"); password != "" {
			cmdArgs = append(cmdArgs, "--password", password)
		}

		// Database
		if database, _ := cmd.Flags().GetString("database"); database != "" {
			cmdArgs = append(cmdArgs, "--database", database)
		} else {
			cmdArgs = append(cmdArgs, "--database", "benchmark")
		}

		// File
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			cmdArgs = append(cmdArgs, "--file", file)
		} else {
			cmdArgs = append(cmdArgs, "--file", filepath.Join(os.TempDir(), "tsbs_data"))
		}

		// Batch size
		if batchSize, _ := cmd.Flags().GetInt("batch-size"); batchSize > 0 {
			cmdArgs = append(cmdArgs, "--batch-size", fmt.Sprintf("%d", batchSize))
		}

		// Workers
		if workers, _ := cmd.Flags().GetInt("workers"); workers > 0 {
			cmdArgs = append(cmdArgs, "--workers", fmt.Sprintf("%d", workers))
		}

		fmt.Printf("Executing: %s\n", strings.Join(cmdArgs, " "))

		execCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Dir = tsbsPath
		if err := execCmd.Run(); err != nil {
			fmt.Printf("Error loading data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nData loading completed!")
	},
}

// TSBSRunQueriesCmd runs query benchmarks
var TSBSRunQueriesCmd = &cobra.Command{
	Use:   "run",
	Short: "Run query benchmarks",
	Long: `Run benchmark queries against KWDB.

Examples:
  # Run all queries from file
  kwcli tsbs run --file=/tmp/tsbs_queries`,
	Run: func(cmd *cobra.Command, args []string) {
		tsbsPath, err := getTSBSPath()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if tsbsPath == "docker" {
			fmt.Println("Docker mode: use 'kwcli tsbs docker-run' instead")
			os.Exit(1)
		}

		binaryPath := filepath.Join(tsbsPath, "tsbs_run_queries_kwdb")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			binaryPath = binaryPath + ".exe"
		}

		// Build command
		cmdArgs := []string{binaryPath}

		// Host
		if host, _ := cmd.Flags().GetString("host"); host != "" {
			cmdArgs = append(cmdArgs, "--host", host)
		}

		// Port
		if port, _ := cmd.Flags().GetInt("port"); port > 0 {
			cmdArgs = append(cmdArgs, "--port", fmt.Sprintf("%d", port))
		}

		// User
		if user, _ := cmd.Flags().GetString("user"); user != "" {
			cmdArgs = append(cmdArgs, "--user", user)
		}

		// Password
		if password, _ := cmd.Flags().GetString("password"); password != "" {
			cmdArgs = append(cmdArgs, "--password", password)
		}

		// Database
		if database, _ := cmd.Flags().GetString("database"); database != "" {
			cmdArgs = append(cmdArgs, "--database", database)
		} else {
			cmdArgs = append(cmdArgs, "--database", "benchmark")
		}

		// File
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			cmdArgs = append(cmdArgs, "--file", file)
		} else {
			cmdArgs = append(cmdArgs, "--file", filepath.Join(os.TempDir(), "tsbs_queries"))
		}

		// Workers
		if workers, _ := cmd.Flags().GetInt("workers"); workers > 0 {
			cmdArgs = append(cmdArgs, "--workers", fmt.Sprintf("%d", workers))
		}

		// Limit
		if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
			cmdArgs = append(cmdArgs, "--limit", fmt.Sprintf("%d", limit))
		}

		fmt.Printf("Executing: %s\n", strings.Join(cmdArgs, " "))

		execCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Dir = tsbsPath
		if err := execCmd.Run(); err != nil {
			fmt.Printf("Error running queries: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nQuery benchmark completed!")
	},
}

// TSBSListCmd lists available query types
var TSBSListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available query types",
	Long:  "List all available TSBS query types for different use cases.",
	Run: func(cmd *cobra.Command, args []string) {
		useCase, _ := cmd.Flags().GetString("use-case")

		fmt.Println("Available TSBS Query Types:")
		fmt.Println(strings.Repeat("-", 60))

		queryTypes := getQueryTypes(useCase)
		for _, qt := range queryTypes {
			fmt.Printf("  %-30s %s\n", qt.Name, qt.Description)
		}

		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  kwcli tsbs generate-queries --use-case=cpu --query-type=double-groupby")
		fmt.Println("  kwcli tsbs generate-queries --use-case=iot --query-type=threshold")
	},
}

func getQueryTypes(useCase string) []struct {
	Name        string
	Description string
} {
	if useCase == "iot" {
		return []struct {
			Name        string
			Description string
		}{
			{"iot-ingest", "Measure data ingestion rate"},
			{"iot-get-last-n", "Get last N readings for a device"},
			{"iot-get-n-historical", "Get N historical readings"},
			{"iot-get-n-moving", "Get N moving average"},
			{"iot-threshold", "Find readings above threshold"},
			{"iot-threshold-agg", "Find aggregates above threshold"},
			{"iot-groupby", "Group by device identifier"},
			{"iot-groupby-time", "Group by time interval"},
			{"iot-double-groupby", "Group by two fields"},
		}
	}

	// Default: cpu use case
	return []struct {
		Name        string
		Description string
	}{
		{"cpu-max-all", "Maximum CPU usage across all devices"},
		{"cpu-max-by-host", "Maximum CPU usage per host"},
		{"double-groupby", "Group by two fields"},
		{"groupby-time", "Group by time interval"},
		{"high-cpu-all", "Find all high CPU readings"},
		{"high-cpu-by-host", "Find high CPU per host"},
		{"lastpoint", "Get last reading per device"},
		{"multi-measure-query-all", "Query multiple measures"},
		{"multi-measure-where", "Query with WHERE clause"},
		{"single-groupby-agg", "Single group with aggregation"},
		{"single-groupby-agg-five", "Five groupings with aggregation"},
		{"single-groupby-agg-max", "Max aggregation per group"},
		{"single-groupby-agg-min", "Min aggregation per group"},
		{"single-groupby-raw", "Raw data per group"},
		{"time-range", "Query time range"},
		{"time-series-all", "All time series data"},
	}
}

func init() {
	rootCmd.AddCommand(tsbsCmd)

	// Add subcommands
	tsbsCmd.AddCommand(TSBSInitCmd)
	tsbsCmd.AddCommand(TSBSLoadCmd)
	tsbsCmd.AddCommand(TSBSRunQueriesCmd)
	tsbsCmd.AddCommand(TSBSListCmd)

	// TSBSInitCmd flags
	TSBSInitCmd.Flags().StringP("use-case", "u", "cpu", "Use case (cpu, iot)")
	TSBSInitCmd.Flags().IntP("scale", "s", 10, "Number of devices")
	TSBSInitCmd.Flags().String("timestamp-start", "2024-01-01T00:00:00Z", "Start timestamp (RFC3339)")
	TSBSInitCmd.Flags().String("timestamp-end", "2024-01-02T00:00:00Z", "End timestamp (RFC3339)")
	TSBSInitCmd.Flags().String("sampling-interval", "10s", "Sampling interval")
	TSBSInitCmd.Flags().Int("queries", 1000, "Number of queries to generate")
	TSBSInitCmd.Flags().String("query-type", "", "Query type (e.g., double-groupby)")
	TSBSInitCmd.Flags().String("data-file", "", "Output data file path")
	TSBSInitCmd.Flags().String("query-file", "", "Output query file path")
	TSBSInitCmd.Flags().String("output-dir", "", "Output directory for both data and query files")
	TSBSInitCmd.Flags().Int("days", 1, "Number of days for data generation")

	TSBSLoadCmd.Flags().String("host", "127.0.0.1", "KWDB host")
	TSBSLoadCmd.Flags().Int("port", 50000, "KWDB port")
	TSBSLoadCmd.Flags().String("user", "root", "KWDB user")
	TSBSLoadCmd.Flags().String("password", "root", "KWDB password")
	TSBSLoadCmd.Flags().String("database", "benchmark", "Database name")
	TSBSLoadCmd.Flags().StringP("file", "f", "", "Input data file")
	TSBSLoadCmd.Flags().Int("batch-size", 10000, "Batch size")
	TSBSLoadCmd.Flags().Int("workers", 4, "Number of workers")

	TSBSRunQueriesCmd.Flags().String("host", "127.0.0.1", "KWDB host")
	TSBSRunQueriesCmd.Flags().Int("port", 50000, "KWDB port")
	TSBSRunQueriesCmd.Flags().String("user", "root", "KWDB user")
	TSBSRunQueriesCmd.Flags().String("password", "root", "KWDB password")
	TSBSRunQueriesCmd.Flags().String("database", "benchmark", "Database name")
	TSBSRunQueriesCmd.Flags().StringP("file", "f", "", "Query file path")
	TSBSRunQueriesCmd.Flags().Int("workers", 4, "Number of workers")
	TSBSRunQueriesCmd.Flags().Int("limit", 0, "Limit number of queries to run")

	TSBSListCmd.Flags().StringP("use-case", "u", "", "Filter by use case (cpu, iot)")
}
