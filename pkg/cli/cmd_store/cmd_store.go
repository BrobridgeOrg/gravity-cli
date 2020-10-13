package cmd_store

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gravitydbs = []string{
	"postgreSQL",
	"mysql-1",
	"mysql-2",
	"MSSQL",
	"BrobridgeSQL",
}

var DbFlag string
var AllFlag bool

var ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all database control by gravity",
	Long:  `List all database control by gravity`,
	Run: func(cmd *cobra.Command, args []string) {
		if AllFlag {

			fmt.Println("=========ALL DATABASE=========")
			for _, db := range gravitydbs {
				fmt.Println(db)
			}
			fmt.Println("==============================")
		} else {
			// TODO: check dbFLag is null

			fmt.Println("=========LIST DATABASE========")
			_, found := Find(gravitydbs, DbFlag)
			if !found {
				fmt.Println("Gravity not manage this db !")
			} else {
				fmt.Println(DbFlag)
			}
			fmt.Println("==============================")
		}
	},
}

var StoreCmd = &cobra.Command{
	Use:   "store",
	Short: "Store function for ls/recover for database",
	Long:  `Store function for ls/recover for database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=============STORE============")
	},
}

var RecoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover all gravity database",
	Long:  `Recover all gravity database`,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if AllFlag {

			fmt.Println("==== RECOVER ALL DATABASE ====")
			for _, db := range gravitydbs {
				fmt.Println(db + " Recovered")
			}
			fmt.Println("==============================")
		} else {

			fmt.Println("=======RECOVER DATABASE=======")
			_, found := Find(gravitydbs, DbFlag)
			if !found {
				fmt.Println("Gravity not manage this db !")
			} else {
				fmt.Println(DbFlag)
				fmt.Println("==============================")

			}
		}
	},
}

func NewStoreCmd() *cobra.Command {

	// Store list Command
	ListCmd.Flags().StringVarP(&DbFlag, "dbinfo", "d", "", "Database information")
	ListCmd.Flags().BoolVarP(&AllFlag, "alldbinfo", "a", false, "All Database information")

	// Store recover Command
	RecoverCmd.Flags().StringVarP(&DbFlag, "dbinfo", "d", "", "Recover Database")
	RecoverCmd.Flags().BoolVarP(&AllFlag, "alldbinfo", "a", false, "Recover all Database")

	// Store Combine
	StoreCmd.AddCommand(ListCmd)
	StoreCmd.AddCommand(RecoverCmd)

	return StoreCmd
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
