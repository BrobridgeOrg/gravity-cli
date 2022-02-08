package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-cli/pkg/logger"
	"github.com/BrobridgeOrg/gravity-cli/pkg/product"
	product_sdk "github.com/BrobridgeOrg/gravity-sdk/product"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type productCmdFunc func(*configs.Config, *zap.Logger, *connector.Connector, *product.Product, *cobra.Command, []string)

// Product flags
var productName string
var productDesc string
var productEnabled bool

// Rule flags
var ruleName string
var ruleDescription string
var ruleEvent string
var ruleMethod string
var ruleEnabled bool
var rulePrimaryKey []string

func init() {

	RootCmd.AddCommand(productCmd)
	productCmd.AddCommand(productListCmd)
	productCmd.AddCommand(productCreateCmd)
	productCmd.AddCommand(productDeleteCmd)

	// Update product
	productCmd.AddCommand(productUpdateCmd)
	productUpdateCmd.Flags().StringVar(&productDesc, "desc", "", "Specify description")
	productUpdateCmd.Flags().BoolVar(&productEnabled, "enabled", false, "Enable product (default false)")

	// Rule
	productCmd.AddCommand(productRuleCmd)

	// List rules
	productRuleCmd.AddCommand(productRuleListCmd)
	productRuleListCmd.Flags().StringVar(&productName, "product", "", "Specify product name (required)")
	productRuleListCmd.MarkFlagRequired("product")

	// Create rule
	productRuleCmd.AddCommand(productRuleCreateCmd)
	productRuleCreateCmd.Flags().StringVar(&productName, "product", "", "Specify product name (required)")
	productRuleCreateCmd.Flags().StringVar(&ruleEvent, "event", "", "Specify event name")
	productRuleCreateCmd.Flags().StringVar(&ruleMethod, "method", "", "Specify method (required)")
	productRuleCreateCmd.Flags().BoolVar(&ruleEnabled, "enabled", false, "Enable rule (default false)")
	productRuleCreateCmd.Flags().StringVar(&ruleDescription, "desc", "", "Specify description")
	productRuleCreateCmd.Flags().StringSliceVar(&rulePrimaryKey, "pk", []string{}, `Specify primary key (support multiple fields with separator ",")`)
	productRuleCreateCmd.MarkFlagRequired("product")
	productRuleCreateCmd.MarkFlagRequired("event")
	productRuleCreateCmd.MarkFlagRequired("method")

	// Update rule
	productRuleCmd.AddCommand(productRuleUpdateCmd)
	productRuleUpdateCmd.Flags().StringVar(&productName, "product", "", "Specify product name (required)")
	productRuleUpdateCmd.Flags().StringVar(&ruleEvent, "event", "", "Specify event name")
	productRuleUpdateCmd.Flags().StringVar(&ruleMethod, "method", "", "Specify method")
	productRuleUpdateCmd.Flags().BoolVar(&ruleEnabled, "enabled", false, "Enable rule (default false)")
	productRuleUpdateCmd.Flags().StringVar(&ruleDescription, "desc", "", "Specify description")
	productRuleUpdateCmd.Flags().StringSliceVar(&rulePrimaryKey, "pk", []string{}, `Specify primary key (support multiple fields with separator ",")`)
	productRuleUpdateCmd.MarkFlagRequired("product")

	// Delete rule
	productRuleCmd.AddCommand(productRuleDeleteCmd)
	productRuleDeleteCmd.Flags().StringVar(&productName, "product", "", "Specify product name (required)")
	productRuleDeleteCmd.MarkFlagRequired("product")
}

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Manage data products",
}

func runProductCmd(fn productCmdFunc, cmd *cobra.Command, args []string) error {

	config.SetHost(host)

	fx.New(
		fx.Supply(config),
		fx.Provide(
			logger.GetLogger,
			connector.New,
			product.New,
		),
		fx.Supply(cmd),
		fx.Supply(args),
		fx.Invoke(fn),
		fx.NopLogger,
	).Run()

	return nil
}

var productListCmd = &cobra.Command{
	Use:   "list",
	Short: "List data products",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductListCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductListCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	products, err := p.GetClient().ListProducts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	if len(products) == 0 {
		fmt.Println("No products")
		os.Exit(0)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Description",
		"Status",
		"Rules",
		"Updated",
		"Created",
	})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, product := range products {
		var status string
		if product.Enabled {
			status = "enabled"
		} else {
			status = "disabled"
		}

		table.Append([]string{
			product.Name,
			product.Description,
			status,
			fmt.Sprintf("%d", len(product.Rules)),
			product.UpdatedAt.String(),
			product.CreatedAt.String(),
		})
	}

	table.Render()

	os.Exit(0)
}

var productCreateCmd = &cobra.Command{
	Use:   "create [product name] [description]",
	Short: "Create a new data product",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductCreateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductCreateCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	setting := product_sdk.ProductSetting{}

	if len(args) > 0 {
		setting.Name = args[0]
	}

	if len(args) == 2 {
		setting.Description = args[1]
	}

	_, err := p.GetClient().CreateProduct(&setting)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Printf("Product \"%s\" was created\n", setting.Name)

	os.Exit(0)
}

var productDeleteCmd = &cobra.Command{
	Use:   "delete [product name]",
	Short: "Delete a data product",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductDeleteCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductDeleteCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		os.Exit(1)
		return
	}

	productName = args[0]

	err := p.GetClient().DeleteProduct(productName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Printf("Product \"%s\" was deleted\n", productName)

	os.Exit(0)
}

var productUpdateCmd = &cobra.Command{
	Use:   "update [product name]",
	Short: "Update a data product",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductUpdateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductUpdateCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	productName = args[0]

	// Getting product information
	product, err := p.GetClient().GetProduct(productName)
	if err != nil {
		fmt.Printf("Not found product \"%s\"\n", productName)
		os.Exit(1)
		return
	}

	// Update description
	if cmd.Flags().Changed("desc") {
		product.Description = productDesc
	}

	// Update enabled
	if cmd.Flags().Changed("enabled") {
		product.Enabled = productEnabled
	}

	// Update
	_, err = p.GetClient().UpdateProduct(productName, product)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	os.Exit(0)
}

var productRuleCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "Manage rules of data product",
}

var productRuleCreateCmd = &cobra.Command{
	Use:   "create [rule name]",
	Short: "Create a new rule",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleCreateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleCreateCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	ruleName = args[0]

	// Validate product name
	if len(productName) == 0 {
		fmt.Println("Require product")
		os.Exit(1)
		return
	}

	if len(ruleEvent) == 0 {
		fmt.Println("Require event")
		os.Exit(1)
		return
	}

	if len(ruleMethod) == 0 {
		fmt.Println("Require method")
		os.Exit(1)
		return
	}

	// Getting product information
	product, err := p.GetClient().GetProduct(productName)
	if err != nil {
		fmt.Printf("Not found product \"%s\"\n", productName)
		os.Exit(1)
		return
	}

	if product.Rules == nil {
		product.Rules = make(map[string]*product_sdk.Rule)
	} else {

		// Check whether rule does exist or not
		_, ok := product.Rules[ruleName]
		if ok {
			fmt.Printf("Rule \"%s\" exists already\n", ruleName)
			os.Exit(1)
			return
		}
	}

	// Preparing a new rule
	rule := product_sdk.NewRule()
	rule.Name = ruleName
	rule.Product = productName
	rule.UpdatedAt = time.Now()
	rule.CreatedAt = time.Now()

	// Unique ID
	id, _ := uuid.NewUUID()
	rule.ID = id.String()
	rule.Event = ruleEvent
	rule.Method = ruleMethod
	rule.PrimaryKey = rulePrimaryKey
	rule.Description = ruleDescription

	// Update enabled
	if cmd.Flags().Changed("enabled") {
		rule.Enabled = ruleEnabled
	}

	// Add to rule set
	product.Rules[rule.Name] = rule

	// Update
	_, err = p.GetClient().UpdateProduct(productName, product)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	os.Exit(0)
}

var productRuleUpdateCmd = &cobra.Command{
	Use:   "update [rule name]",
	Short: "Update rule",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleUpdateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleUpdateCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	ruleName = args[0]

	// Validate product name
	if len(productName) == 0 {
		fmt.Println("Require product")
		os.Exit(1)
		return
	}

	// Getting product information
	product, err := p.GetClient().GetProduct(productName)
	if err != nil {
		fmt.Printf("Not found product \"%s\"\n", productName)
		os.Exit(1)
		return
	}

	if product.Rules == nil {
		fmt.Printf("Not found rule \"%s\"\n", ruleName)
		os.Exit(1)
		return
	}

	// Check whether rule does exist or not
	rule, ok := product.Rules[ruleName]
	if !ok {
		fmt.Printf("Not found rule \"%s\"\n", ruleName)
		os.Exit(1)
		return
	}

	if cmd.Flags().Changed("event") {

		if len(ruleEvent) == 0 {
			fmt.Println("Invalid event")
			os.Exit(1)
			return
		}

		rule.Event = ruleEvent
	}

	if cmd.Flags().Changed("method") {

		if len(ruleEvent) == 0 {
			fmt.Println("Invalid method")
			os.Exit(1)
			return
		}

		rule.Method = ruleMethod
	}

	if cmd.Flags().Changed("pk") {
		rule.PrimaryKey = rulePrimaryKey
	}

	if cmd.Flags().Changed("desc") {
		rule.Description = ruleDescription
	}

	// Update enabled
	if cmd.Flags().Changed("enabled") {
		rule.Enabled = ruleEnabled
	}

	rule.UpdatedAt = time.Now()

	// Update
	_, err = p.GetClient().UpdateProduct(productName, product)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	os.Exit(0)
}

var productRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List rules of product",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleListCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleListCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	// Validate product name
	if len(productName) == 0 {
		fmt.Println("Require product")
		os.Exit(1)
		return
	}

	// Getting product information
	product, err := p.GetClient().GetProduct(productName)
	if err != nil {
		fmt.Printf("Not found product \"%s\"\n", productName)
		os.Exit(1)
		return
	}

	if product.Rules == nil {
		product.Rules = make(map[string]*product_sdk.Rule)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Description",
		"Event",
		"Method",
		"PK",
		"Status",
		"Updated",
		"Created",
	})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, rule := range product.Rules {
		var pk string
		if len(rule.PrimaryKey) == 0 {
			pk = "n/a"
		} else {
			pk = strings.Join(rule.PrimaryKey, ",")
		}

		var status string
		if rule.Enabled {
			status = "enabled"
		} else {
			status = "disabled"
		}

		table.Append([]string{
			rule.Name,
			rule.Description,
			rule.Event,
			rule.Method,
			pk,
			status,
			rule.UpdatedAt.String(),
			rule.CreatedAt.String(),
		})
	}

	table.Render()

	os.Exit(0)
}

var productRuleDeleteCmd = &cobra.Command{
	Use:   "delete [rule name]",
	Short: "Delete rule",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleDeleteCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleDeleteCmd(config *configs.Config, l *zap.Logger, c *connector.Connector, p *product.Product, cmd *cobra.Command, args []string) {

	ruleName = args[0]

	// Validate product name
	if len(productName) == 0 {
		fmt.Println("Require product")
		os.Exit(1)
		return
	}

	// Getting product information
	product, err := p.GetClient().GetProduct(productName)
	if err != nil {
		fmt.Printf("Not found product \"%s\"\n", productName)
		os.Exit(1)
		return
	}

	if product.Rules == nil {
		fmt.Printf("Not found rule \"%s\"\n", ruleName)
		os.Exit(1)
		return
	}

	// Check whether rule does exist or not
	_, ok := product.Rules[ruleName]
	if !ok {
		fmt.Printf("Not found rule \"%s\"\n", ruleName)
		os.Exit(1)
		return
	}

	delete(product.Rules, ruleName)

	// Update
	_, err = p.GetClient().UpdateProduct(productName, product)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	os.Exit(0)
}
