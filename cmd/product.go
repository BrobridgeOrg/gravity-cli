package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	record_type "github.com/BrobridgeOrg/compton/types/record"
	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-cli/pkg/logger"
	"github.com/BrobridgeOrg/gravity-cli/pkg/product"
	product_sdk "github.com/BrobridgeOrg/gravity-sdk/product"
	subscriber_sdk "github.com/BrobridgeOrg/gravity-sdk/subscriber"
	gravity_sdk_types_product_event "github.com/BrobridgeOrg/gravity-sdk/types/product_event"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	productEventStream = "GVT_%s_DP_%s"
)

type ProductCommandContext struct {
	Config    *configs.Config
	Logger    *zap.Logger
	Connector *connector.Connector
	Product   *product.Product
	Cmd       *cobra.Command
	Args      []string
}

type productCmdFunc func(*ProductCommandContext) error

// Product flags
var productName string
var productDesc string
var productEnabled bool
var productSchemaFile string

// product subscriber
var productSubscriberName string
var productSubscriberStartSeq uint64
var productSubscriberPartitions []int

// Rule flags
var ruleName string
var ruleDescription string
var ruleEvent string
var ruleMethod string
var ruleEnabled bool
var rulePrimaryKey []string
var ruleSchemaFile string
var ruleHandlerFile string

func init() {

	RootCmd.AddCommand(productCmd)
	productCmd.AddCommand(productListCmd)

	// Create product
	productCmd.AddCommand(productCreateCmd)
	productCreateCmd.Flags().StringVar(&productDesc, "desc", "", "Specify description")
	productCreateCmd.Flags().BoolVar(&productEnabled, "enabled", false, "Enable product (default false)")
	productCreateCmd.Flags().StringVar(&productSchemaFile, "schema", "", "Load schema from specific file")

	// Update product
	productCmd.AddCommand(productUpdateCmd)
	productUpdateCmd.Flags().StringVar(&productDesc, "desc", "", "Specify description")
	productUpdateCmd.Flags().BoolVar(&productEnabled, "enabled", false, "Enable produc")
	productUpdateCmd.Flags().StringVar(&productSchemaFile, "schema", "", "Load schema from specific file")

	// Delete and purge product
	productCmd.AddCommand(productDeleteCmd)
	productCmd.AddCommand(productPurgeCmd)

	// Show product information
	productCmd.AddCommand(productInfoCmd)

	// Subscribe product
	productCmd.AddCommand(productSubCmd)
	productSubCmd.Flags().StringVar(&productSubscriberName, "name", "", "Specify subscriber name")
	productSubCmd.Flags().Uint64Var(&productSubscriberStartSeq, "seq", 1, "Specify start sequence")
	productSubCmd.Flags().IntSliceVar(&productSubscriberPartitions, "partitions", []int{-1}, "Specify partitions (default -1 for all)")

	// Snapshot
	productCmd.AddCommand(productSnapshotCmd)

	// Rule
	productCmd.AddCommand(productRuleCmd)

	// List rules
	productRuleCmd.AddCommand(productRuleListCmd)

	// Create rule
	productRuleCmd.AddCommand(productRuleAddCmd)
	productRuleAddCmd.Flags().StringVar(&ruleEvent, "event", "", "Specify event name")
	productRuleAddCmd.Flags().StringVar(&ruleMethod, "method", "", "Specify method (required)")
	productRuleAddCmd.Flags().BoolVar(&ruleEnabled, "enabled", false, "Enable rule (default false)")
	productRuleAddCmd.Flags().StringVar(&ruleDescription, "desc", "", "Specify description")
	productRuleAddCmd.Flags().StringSliceVar(&rulePrimaryKey, "pk", []string{}, `Specify primary key (support multiple fields with separator ",")`)
	productRuleAddCmd.Flags().StringVar(&ruleSchemaFile, "schema", "", "Load schema from specific file")
	productRuleAddCmd.Flags().StringVar(&ruleHandlerFile, "handler", "", "Load handler script from specific file")
	productRuleAddCmd.MarkFlagRequired("event")
	productRuleAddCmd.MarkFlagRequired("method")

	// Update rule
	productRuleCmd.AddCommand(productRuleUpdateCmd)
	productRuleUpdateCmd.Flags().StringVar(&ruleEvent, "event", "", "Specify event name")
	productRuleUpdateCmd.Flags().StringVar(&ruleMethod, "method", "", "Specify method")
	productRuleUpdateCmd.Flags().BoolVar(&ruleEnabled, "enabled", false, "Enable rule (default false)")
	productRuleUpdateCmd.Flags().StringVar(&ruleDescription, "desc", "", "Specify description")
	productRuleUpdateCmd.Flags().StringSliceVar(&rulePrimaryKey, "pk", []string{}, `Specify primary key (support multiple fields with separator ",")`)
	productRuleUpdateCmd.Flags().StringVar(&ruleSchemaFile, "schema", "", "Load schema from specific file")
	productRuleUpdateCmd.Flags().StringVar(&ruleHandlerFile, "handler", "", "Load handler script from specific file")

	// Delete rule
	productRuleCmd.AddCommand(productRuleDeleteCmd)

	// Show rule information
	productRuleCmd.AddCommand(productRuleInfoCmd)
}

func readSchemaFile(filename string) (map[string]interface{}, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("No such schema file")
	}

	// Read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	err = json.Unmarshal(data, &schema)
	if err != nil {

		return nil, errors.New("invalid schema format")
	}

	return schema, nil
}

func readHandlerScriptFile(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("No such handler file")
	}

	// Read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Manage data products",
}

func runProductCmd(fn productCmdFunc, cmd *cobra.Command, args []string) error {

	var cctx *ProductCommandContext

	config.SetHost(host)

	app := fx.New(
		fx.Supply(config),
		fx.Provide(
			logger.GetLogger,
			connector.New,
			product.New,
		),
		fx.Supply(cmd),
		fx.Supply(args),
		fx.Provide(func(
			config *configs.Config,
			l *zap.Logger,
			c *connector.Connector,
			p *product.Product,
			cmd *cobra.Command,
			args []string,
		) *ProductCommandContext {
			return &ProductCommandContext{
				Config:    config,
				Logger:    l,
				Connector: c,
				Product:   p,
				Cmd:       cmd,
				Args:      args,
			}
		}),
		fx.Populate(&cctx),
		fx.NopLogger,
	)

	if err := app.Start(context.Background()); err != nil {
		return err
	}

	return fn(cctx)
}

var productListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available products",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductListCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductListCmd(cctx *ProductCommandContext) error {

	products, err := cctx.Product.GetClient().ListProducts()
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	if len(products) == 0 {
		cctx.Cmd.SilenceUsage = true
		return errors.New("No available products")
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

	return nil
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

func runProductCreateCmd(cctx *ProductCommandContext) error {

	setting := product_sdk.ProductSetting{}
	setting.Name = cctx.Args[0]

	// Description
	if cctx.Cmd.Flags().Changed("desc") {
		setting.Description = productDesc
	}

	// Enable
	if cctx.Cmd.Flags().Changed("enabled") {
		setting.Enabled = productEnabled
	}

	// Schema
	if cctx.Cmd.Flags().Changed("schema") {
		schema, err := readSchemaFile(productSchemaFile)
		if err != nil {
			return err
		}

		setting.Schema = schema
	}

	setting.Stream = fmt.Sprintf(productEventStream, domain, setting.Name)

	// Snapshot
	setting.EnabledSnapshot = true

	_, err := cctx.Product.GetClient().CreateProduct(&setting)
	if err != nil {
		return err
	}

	fmt.Printf("Product \"%s\" was created\n", setting.Name)

	return nil
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

func runProductDeleteCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]

	err := cctx.Product.GetClient().DeleteProduct(productName)
	if err != nil {
		return err
	}

	fmt.Printf("Product \"%s\" was deleted\n", productName)

	return nil
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

func runProductUpdateCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]
	changed := false

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	// Update description
	if cctx.Cmd.Flags().Changed("desc") {
		product.Description = productDesc
		changed = true
	}

	// Update enabled
	if cctx.Cmd.Flags().Changed("enabled") {
		product.Enabled = productEnabled
		changed = true
	}

	// Update schema
	if cctx.Cmd.Flags().Changed("schema") {
		schema, err := readSchemaFile(productSchemaFile)
		if err != nil {
			return err
		}

		product.Schema = schema
		changed = true
	}

	// Nothing's changed
	if !changed {
		return nil
	}

	product.Stream = fmt.Sprintf(productEventStream, domain, productName)

	product.EnabledSnapshot = true

	// Update
	_, err = cctx.Product.GetClient().UpdateProduct(productName, product)
	if err != nil {
		return err
	}

	return nil
}

var productPurgeCmd = &cobra.Command{
	Use:   "purge [product name]",
	Short: "purge a data product without deleting it",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductPurgeCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductPurgeCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]

	err := cctx.Product.GetClient().PurgeProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	fmt.Printf("Product \"%s\" was purged\n", productName)

	return nil
}

var productInfoCmd = &cobra.Command{
	Use:   "info [product name]",
	Short: "Show information about product",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductInfoCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductInfoCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	//	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	//table.SetNoWhiteSpace(true)

	// Basic information
	table.Append([]string{
		"Name:",
		product.Name,
	})

	table.Append([]string{
		"Description:",
		product.Description,
	})

	var status string
	if product.Enabled {
		status = "enabled"
	} else {
		status = "disabled"
	}

	table.Append([]string{
		"Status:",
		status,
	})

	table.Append([]string{
		"Stream:",
		product.Stream,
	})

	table.Append([]string{
		"Updated:",
		product.UpdatedAt.String(),
	})

	table.Append([]string{
		"Created:",
		product.CreatedAt.String(),
	})

	fmt.Printf("Information for Product %s\n\n", productName)
	fmt.Printf("Configuration:\n\n")

	table.Render()

	fmt.Println("")

	// Schema
	if product.Schema != nil {
		fmt.Println("")
		fmt.Println("Schema:")
		fmt.Println("")
		schema, _ := json.MarshalIndent(product.Schema, "", "    ")
		fmt.Println(string(schema))
		fmt.Println("")
	}

	return nil
}

var productSubCmd = &cobra.Command{
	Use:   "sub [product name]",
	Short: "Generic subscription client for product",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductSubCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductSubCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]

	fmt.Printf("Subscribing to product: %s\n", productName)

	if len(productSubscriberName) > 0 {
		fmt.Printf("Subscriber name: %s\n", productSubscriberName)
	}

	fmt.Printf("Start Sequence: %d\n", productSubscriberStartSeq)

	// Initializing gravity subscriber
	opts := subscriber_sdk.NewOptions()
	opts.Verbose = true
	s := subscriber_sdk.NewSubscriberWithClient(productSubscriberName, cctx.Connector.GetClient(), opts)
	_, err := s.Subscribe(productName, func(msg *nats.Msg) {

		var pe gravity_sdk_types_product_event.ProductEvent

		err := proto.Unmarshal(msg.Data, &pe)
		if err != nil {
			fmt.Printf("Failed to parsing product event: %v", err)
			msg.Ack()
			return
		}

		md, _ := msg.Metadata()

		r, err := pe.GetContent()
		if err != nil {
			fmt.Printf("Failed to parsing content: %v", err)
			msg.Ack()
			return
		}

		// Convert data to JSON
		event := map[string]interface{}{
			"header":     msg.Header,
			"subject":    msg.Subject,
			"seq":        md.Sequence.Consumer,
			"timestamp":  md.Timestamp,
			"product":    productName,
			"event":      pe.EventName,
			"method":     pe.Method.String(),
			"table":      pe.Table,
			"primaryKey": pe.PrimaryKeys,
			"payload":    r.AsMap(),
		}

		data, _ := json.MarshalIndent(event, "", "  ")
		fmt.Println(string(data))
	}, subscriber_sdk.Partition(productSubscriberPartitions...), subscriber_sdk.StartSequence(productSubscriberStartSeq))
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	select {}

	return nil
}

var productRuleCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "Manage rules of data product",
}

var productRuleAddCmd = &cobra.Command{
	Use:   "add [product] [rule name]",
	Short: "Add a new rule to product",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleAddCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleAddCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]
	ruleName = cctx.Args[1]

	// Validate product name
	if len(productName) == 0 {
		return errors.New(fmt.Sprintf("require product"))
	}

	if len(ruleEvent) == 0 {
		return errors.New(fmt.Sprintf("require flag: --event"))
	}

	if len(ruleMethod) == 0 {
		return errors.New(fmt.Sprintf("require flag: --method"))
	}

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	if product.Rules == nil {
		product.Rules = make(map[string]*product_sdk.Rule)
	} else {

		// Check whether rule does exist or not
		_, ok := product.Rules[ruleName]
		if ok {
			cctx.Cmd.SilenceUsage = true
			return errors.New(fmt.Sprintf("Rule \"%s\" exists already\n", ruleName))
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

	// Enabled
	if cctx.Cmd.Flags().Changed("enabled") {
		rule.Enabled = ruleEnabled
	}

	// Schema
	if cctx.Cmd.Flags().Changed("schema") {
		schema, err := readSchemaFile(ruleSchemaFile)
		if err != nil {
			return err
		}

		rule.SchemaConfig = schema
	}

	// Handler script
	if cctx.Cmd.Flags().Changed("handler") {

		script, err := readHandlerScriptFile(ruleHandlerFile)
		if err != nil {
			return err
		}

		rule.HandlerConfig = &product_sdk.HandlerConfig{
			Type:   "script",
			Script: string(script),
		}
	}

	// Add to rule set
	product.Rules[rule.Name] = rule

	// Update
	_, err = cctx.Product.GetClient().UpdateProduct(productName, product)
	if err != nil {
		return err
	}

	return nil
}

var productRuleUpdateCmd = &cobra.Command{
	Use:   "update [product name] [rule name]",
	Short: "Update rule of product",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleUpdateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleUpdateCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]
	ruleName = cctx.Args[1]

	// Validate product name
	if len(productName) == 0 {
		return errors.New(fmt.Sprintf("require product"))
	}

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	if product.Rules == nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found rule \"%s\"\n", ruleName))
	}

	// Check whether rule does exist or not
	rule, ok := product.Rules[ruleName]
	if !ok {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found rule \"%s\"\n", ruleName))
	}

	if cctx.Cmd.Flags().Changed("event") {

		if len(ruleEvent) == 0 {
			return errors.New("Invalid event")
		}

		rule.Event = ruleEvent
	}

	if cctx.Cmd.Flags().Changed("method") {

		if len(ruleEvent) == 0 {
			return errors.New("Invalid method")
		}

		rule.Method = ruleMethod
	}

	if cctx.Cmd.Flags().Changed("pk") {
		rule.PrimaryKey = rulePrimaryKey
	}

	if cctx.Cmd.Flags().Changed("desc") {
		rule.Description = ruleDescription
	}

	// Update enabled
	if cctx.Cmd.Flags().Changed("enabled") {
		rule.Enabled = ruleEnabled
	}

	// Update schema
	if cctx.Cmd.Flags().Changed("schema") {
		schema, err := readSchemaFile(ruleSchemaFile)
		if err != nil {
			return err
		}

		rule.SchemaConfig = schema
	}

	// Handler script
	if cctx.Cmd.Flags().Changed("handler") {

		script, err := readHandlerScriptFile(ruleHandlerFile)
		if err != nil {
			return err
		}

		rule.HandlerConfig = &product_sdk.HandlerConfig{
			Type:   "script",
			Script: string(script),
		}
	}

	rule.UpdatedAt = time.Now()

	// Update
	_, err = cctx.Product.GetClient().UpdateProduct(productName, product)
	if err != nil {
		return err
	}

	return nil
}

var productRuleListCmd = &cobra.Command{
	Use:   "list [product name]",
	Short: "List available rules",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleListCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleListCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]

	// Validate product name
	if len(productName) == 0 {
		return errors.New(fmt.Sprintf("require product"))
	}

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	if product.Rules == nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New("No available rules")
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

	return nil
}

var productRuleDeleteCmd = &cobra.Command{
	Use:   "delete [product name] [rule name]",
	Short: "Delete rule",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleDeleteCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleDeleteCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]
	ruleName = cctx.Args[1]

	// Validate product name
	if len(productName) == 0 {
		return errors.New(fmt.Sprintf("require product"))
	}

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	if product.Rules == nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found rule \"%s\"\n", ruleName))
	}

	// Check whether rule does exist or not
	_, ok := product.Rules[ruleName]
	if !ok {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found rule \"%s\"\n", ruleName))
	}

	delete(product.Rules, ruleName)

	// Update
	_, err = cctx.Product.GetClient().UpdateProduct(productName, product)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	return nil
}

var productRuleInfoCmd = &cobra.Command{
	Use:   "info [product name] [rule name]",
	Short: "Show information about rule",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductRuleInfoCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductRuleInfoCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]
	ruleName = cctx.Args[1]

	// Validate product name
	if len(productName) == 0 {
		return errors.New(fmt.Sprintf("require product"))
	}

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	if product.Rules == nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found rule \"%s\"\n", ruleName))
	}

	// Check whether rule does exist or not
	rule, ok := product.Rules[ruleName]
	if !ok {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found rule \"%s\"\n", ruleName))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	//	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	//table.SetNoWhiteSpace(true)

	// Basic information
	table.Append([]string{
		"Product:",
		rule.Product,
	})

	table.Append([]string{
		"Name:",
		rule.Name,
	})

	table.Append([]string{
		"Description:",
		rule.Description,
	})

	table.Append([]string{
		"Event:",
		rule.Event,
	})

	var status string
	if rule.Enabled {
		status = "enabled"
	} else {
		status = "disabled"
	}

	table.Append([]string{
		"Status:",
		status,
	})

	table.Append([]string{
		"Method:",
		rule.Method,
	})

	// Primary Key
	var pk string
	if len(rule.PrimaryKey) == 0 {
		pk = "n/a"
	} else {
		pk = strings.Join(rule.PrimaryKey, ",")
	}
	table.Append([]string{
		"Primary Key:",
		pk,
	})

	table.Append([]string{
		"Updated:",
		rule.UpdatedAt.String(),
	})

	table.Append([]string{
		"Created:",
		rule.CreatedAt.String(),
	})

	fmt.Printf("Information for Rule %s\n\n", productName)
	fmt.Printf("Configuration:\n\n")

	table.Render()

	fmt.Println("")

	// Schema
	if rule.SchemaConfig != nil {
		fmt.Println("")
		fmt.Println("Schema:")
		fmt.Println("")
		schema, _ := json.MarshalIndent(rule.SchemaConfig, "", "    ")
		fmt.Println(string(schema))
		fmt.Println("")
	}

	// Handler
	if rule.HandlerConfig != nil {
		fmt.Println("")
		fmt.Println("Handler Script:")
		fmt.Println("")
		fmt.Println(rule.HandlerConfig.Script)
		fmt.Println("")
	}

	return nil
}

var productSnapshotCmd = &cobra.Command{
	Use:   "snapshot [product name]",
	Short: "Take snapshot from product",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runProductSnapshotCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runProductSnapshotCmd(cctx *ProductCommandContext) error {

	productName = cctx.Args[0]

	// Getting product information
	product, err := cctx.Product.GetClient().GetProduct(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found product \"%s\"\n", productName))
	}

	if !product.EnabledSnapshot {
		cctx.Cmd.SilenceUsage = true
		return errors.New("Product snapshot is not enabled")
	}

	// Create new snapshot
	s, err := cctx.Product.GetClient().CreateSnapshot(productName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	// Fetch records
	msgs, err := s.Fetch()
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	for msg := range msgs {

		fmt.Println(string(msg.Data))
		var r record_type.Record
		err := record_type.Unmarshal(msg.Data, &r)
		if err != nil {
			cctx.Cmd.SilenceUsage = true
			return err
		}

		// Convert data to JSON
		data, _ := json.MarshalIndent(r.AsMap(), "", "  ")
		fmt.Println(string(data))

		msg.Ack()
	}

	err = s.Close()
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	return nil
}
