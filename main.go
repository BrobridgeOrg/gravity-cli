package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gravity-cli",
		Short: "Gravity-cli tool",
		Long:  `Gravity-cli tool by daginwu in Brobridge`,
	}

	database string
)

func main() {
	Execute()
}

func init() {

	rootCmd.AddCommand(storeCmd)
	storeCmd.AddCommand(listCmd)
	storeCmd.AddCommand(recoverCmd)

	// listCmd.Flags().StringVarP(&database, "database", "db", "", "Database name")
	// recoverCmd.Flags().StringVarP(&database, "database", "db", "", "Database name")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all database control by gravity",
	Long:  `List all database control by gravity`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==============================")
		fmt.Println("=========LIST DATABASE========")
		fmt.Println("==============================")
	},
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store function for ls/recover for database",
	Long:  `Store function for ls/recover for database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==============================")
		fmt.Println("=============STORE============")
		fmt.Println("==============================")
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all gravity database",
	Long:  `List all gravity database`,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==============================")
		fmt.Println("=========LIST DATABASE========")
		fmt.Println("==============================")
	},
}

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover all gravity database",
	Long:  `Recover all gravity database`,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==============================")
		fmt.Println("=======RECOVER DATABASE=======")
		fmt.Println("==============================")
	},
}
