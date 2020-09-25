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

	gravitydbs = []string{
		"postgreSQL",
		"mysql-1",
		"mysql-2",
		"MSSQL",
		"BrobridgeSQL",
	}

	dbFlag  string
	allFlag bool
)

func main() {
	Execute()
}

func init() {
	listCmd.Flags().StringVarP(&dbFlag, "dbinfo", "d", "", "Database information")
	listCmd.Flags().BoolVarP(&allFlag, "alldbinfo", "a", false, "All Database information")

	recoverCmd.Flags().StringVarP(&dbFlag, "dbinfo", "d", "", "Recover Database")
	recoverCmd.Flags().BoolVarP(&allFlag, "alldbinfo", "a", false, "Recover all Database")

	rootCmd.AddCommand(storeCmd)
	storeCmd.AddCommand(listCmd)
	storeCmd.AddCommand(recoverCmd)

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
		if allFlag {

			fmt.Println("=========ALL DATABASE=========")
			for _, db := range gravitydbs {
				fmt.Println(db)
			}
			fmt.Println("==============================")
		} else {
			// TODO: check dbFLag is null

			fmt.Println("=========LIST DATABASE========")
			_, found := Find(gravitydbs, dbFlag)
			if !found {
				fmt.Println("Gravity not manage this db !")
			} else {
				fmt.Println(dbFlag)
			}
			fmt.Println("==============================")
		}
	},
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store function for ls/recover for database",
	Long:  `Store function for ls/recover for database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=============STORE============")
	},
}

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover all gravity database",
	Long:  `Recover all gravity database`,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if allFlag {

			fmt.Println("==== RECOVER ALL DATABASE ====")
			for _, db := range gravitydbs {
				fmt.Println(db + " Recovered")
			}
			fmt.Println("==============================")
		} else {

			fmt.Println("=======RECOVER DATABASE=======")
			_, found := Find(gravitydbs, dbFlag)
			if !found {
				fmt.Println("Gravity not manage this db !")
			} else {
				fmt.Println(dbFlag)
				fmt.Println("==============================")

			}
		}
	},
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
