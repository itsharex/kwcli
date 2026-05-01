package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shawn0915/kwcli/pkg/sampledb"
	"github.com/spf13/cobra"
)

// sampledbCmd represents the sampledb command
var sampledbCmd = &cobra.Command{
	Use:   "sampledb",
	Short: "Manage KWDB SampleDB (Smart Meter model)",
	Long:  "Initialize schema, generate data, and run scenario queries for the KWDB Smart Meter sample database.",
}

var sampledbInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create databases and tables for Smart Meter model",
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := sampledb.NewRunner()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Creating databases and tables...")
		fmt.Println("  [1/2] Creating relational database (rdb)...")
		if err := runner.ExecBatch(sampledb.RDBSchema); err != nil {
			fmt.Printf("Failed to create RDB schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("  [2/2] Creating time-series database (tsdb)...")
		if err := runner.ExecBatch(sampledb.TSDBSchema); err != nil {
			fmt.Printf("Failed to create TSDB schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Schema initialized successfully!")
		fmt.Println("Run 'kwcli sampledb generate' to populate sample data.")
	},
}

var sampledbGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate sample data for Smart Meter model",
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := sampledb.NewRunner()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Generating sample data...")
		fmt.Println("  [1/2] Generating relationship data (users, areas, meters, rules)...")
		if err := runner.ExecBatch(sampledb.GenerateRDBData()); err != nil {
			fmt.Printf("Failed to generate RDB data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("  [2/2] Generating time-series data (10,000 readings)...")
		if err := runner.ExecBatch(sampledb.GenerateTSDBData()); err != nil {
			fmt.Printf("Failed to generate TSDB data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Sample data generated successfully!")
		fmt.Println("Run 'kwcli sampledb list' to see available scenarios.")
	},
}

var sampledbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available scenario queries",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available Smart Meter Scenarios:")
		fmt.Println(strings.Repeat("-", 60))
		for i, s := range sampledb.AllScenarios {
			fmt.Printf("  %2d. %-25s %s\n", i+1, s.Name, s.Title)
			fmt.Printf("      %s\n", s.Description)
		}
		fmt.Println()
		fmt.Println("Run a scenario:")
		fmt.Println("  kwcli sampledb run <name>")
		fmt.Println("  kwcli sampledb run --all")
	},
}

var sampledbRunAll bool

var sampledbRunCmd = &cobra.Command{
	Use:   "run <scenario-name>",
	Short: "Run a scenario query",
	Long:  "Run a specific scenario query by name, or use --all to run all scenarios.",
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := sampledb.NewRunner()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if sampledbRunAll {
			fmt.Println("Running all scenarios...")
			fmt.Println()
			for _, s := range sampledb.AllScenarios {
				fmt.Printf("=== %s ===\n", s.Title)
				fmt.Printf("-- %s\n", s.Description)
				fmt.Println()
				if err := runner.Exec(s.SQL); err != nil {
					fmt.Printf("Failed to run scenario '%s': %v\n", s.Name, err)
				}
				fmt.Println()
			}
			return
		}

		if len(args) == 0 {
			fmt.Println("Please specify a scenario name or use --all.")
			fmt.Println("Run 'kwcli sampledb list' to see available scenarios.")
			os.Exit(1)
		}

		name := args[0]
		scenario := sampledb.FindScenario(name)
		if scenario == nil {
			fmt.Printf("Unknown scenario: %s\n", name)
			fmt.Println("Run 'kwcli sampledb list' to see available scenarios.")
			os.Exit(1)
		}

		fmt.Printf("=== %s ===\n", scenario.Title)
		fmt.Printf("-- %s\n", scenario.Description)
		fmt.Println()
		if err := runner.Exec(scenario.SQL); err != nil {
			fmt.Printf("Failed to run scenario: %v\n", err)
			os.Exit(1)
		}
	},
}

var sampledbCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all SampleDB data (drop databases)",
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := sampledb.NewRunner()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Cleaning SampleDB data...")
		fmt.Println("  Dropping tsdb database...")
		if err := runner.Exec("DROP DATABASE IF EXISTS tsdb CASCADE;"); err != nil {
			fmt.Printf("Warning: failed to drop tsdb: %v\n", err)
		}
		fmt.Println("  Dropping rdb database...")
		if err := runner.Exec("DROP DATABASE IF EXISTS rdb CASCADE;"); err != nil {
			fmt.Printf("Warning: failed to drop rdb: %v\n", err)
		}
		fmt.Println("SampleDB cleaned successfully!")
	},
}

var sampledbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if SampleDB schema and data exist",
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := sampledb.NewRunner()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		rdbExists := runner.CheckDatabaseExists("rdb")
		tsdbExists := runner.CheckDatabaseExists("tsdb")

		fmt.Println("SampleDB Status:")
		fmt.Printf("  rdb  database: %s\n", boolStr(rdbExists, "exists", "not found"))
		fmt.Printf("  tsdb database: %s\n", boolStr(tsdbExists, "exists", "not found"))

		if !rdbExists || !tsdbExists {
			fmt.Println()
			fmt.Println("SampleDB is not initialized. Run:")
			fmt.Println("  kwcli sampledb init")
			fmt.Println("  kwcli sampledb generate")
		}
	},
}

func boolStr(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}

func init() {
	rootCmd.AddCommand(sampledbCmd)
	sampledbCmd.AddCommand(sampledbInitCmd)
	sampledbCmd.AddCommand(sampledbGenerateCmd)
	sampledbCmd.AddCommand(sampledbListCmd)
	sampledbCmd.AddCommand(sampledbRunCmd)
	sampledbCmd.AddCommand(sampledbCleanCmd)
	sampledbCmd.AddCommand(sampledbStatusCmd)

	sampledbRunCmd.Flags().BoolVar(&sampledbRunAll, "all", false, "Run all scenarios")
}
