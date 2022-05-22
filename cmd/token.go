package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-cli/pkg/logger"
	"github.com/BrobridgeOrg/gravity-cli/pkg/token"
	token_sdk "github.com/BrobridgeOrg/gravity-sdk/token"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var tokenDesc string
var tokenEnabled bool

type TokenCommandContext struct {
	Config    *configs.Config
	Logger    *zap.Logger
	Connector *connector.Connector
	Token     *token.Token
	Cmd       *cobra.Command
	Args      []string
}

type tokenCmdFunc func(*TokenCommandContext) error

func init() {

	RootCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(tokenListAvailablePermissionsCmd)
	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenDeleteCmd)
	tokenCmd.AddCommand(tokenInfoCmd)

	// Create
	tokenCmd.AddCommand(tokenCreateCmd)
	tokenCreateCmd.Flags().StringVar(&tokenDesc, "desc", "", "Specify description")
	tokenCreateCmd.Flags().BoolVar(&tokenEnabled, "enabled", true, "Enable token (default true)")

	// Update
	tokenCmd.AddCommand(tokenUpdateCmd)
	tokenUpdateCmd.Flags().StringVar(&tokenDesc, "desc", "", "Specify description")
	tokenUpdateCmd.Flags().BoolVar(&tokenEnabled, "enabled", true, "Enable token (default true)")

	// Grant
	tokenCmd.AddCommand(tokenGrantCmd)

	// Revoke
	tokenCmd.AddCommand(tokenRevokeCmd)
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage access tokens",
}

func runTokenCmd(fn tokenCmdFunc, cmd *cobra.Command, args []string) error {

	var cctx *TokenCommandContext

	config.SetHost(host)

	app := fx.New(
		fx.Supply(config),
		fx.Provide(
			logger.GetLogger,
			connector.New,
			token.New,
		),
		fx.Supply(cmd),
		fx.Supply(args),
		fx.Provide(func(
			config *configs.Config,
			l *zap.Logger,
			c *connector.Connector,
			t *token.Token,
			cmd *cobra.Command,
			args []string,
		) *TokenCommandContext {
			return &TokenCommandContext{
				Config:    config,
				Logger:    l,
				Connector: c,
				Token:     t,
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

var tokenListAvailablePermissionsCmd = &cobra.Command{
	Use:   "list_available_permissions",
	Short: "List available permissions",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenListAvailablePermissionsCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenListAvailablePermissionsCmd(cctx *TokenCommandContext) error {

	permissions, err := cctx.Token.GetClient().ListAvailablePermissions()
	if err != nil {
		return err
	}

	if permissions == nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New("No available permissions")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Permissions",
		"Description",
	})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for perm, desc := range permissions {

		table.Append([]string{
			perm,
			desc,
		})
	}

	table.Render()

	return nil
}

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available tokens",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenListCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenListCmd(cctx *TokenCommandContext) error {

	tokens, err := cctx.Token.GetClient().ListTokens()
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		cctx.Cmd.SilenceUsage = true
		return errors.New("No available tokens")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Token ID",
		"Description",
		"Status",
		"Permissions",
		"Updated",
		"Created",
	})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, t := range tokens {
		var status string
		if t.Enabled {
			status = "enabled"
		} else {
			status = "disabled"
		}

		table.Append([]string{
			t.ID,
			t.Description,
			status,
			strconv.Itoa(len(t.Permissions)),
			t.UpdatedAt.String(),
			t.CreatedAt.String(),
		})
	}

	table.Render()

	return nil
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create [description]",
	Short: "Create a new access token",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenCreateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenCreateCmd(cctx *TokenCommandContext) error {

	setting := token_sdk.TokenSetting{}

	// Description
	if cctx.Cmd.Flags().Changed("desc") {
		setting.Description = tokenDesc
	}

	// Enable
	if cctx.Cmd.Flags().Changed("enabled") {
		setting.Enabled = tokenEnabled
	}

	// Generate token ID
	id, _ := uuid.NewUUID()

	token, _, err := cctx.Token.GetClient().CreateToken(id.String(), &setting)
	if err != nil {
		return err
	}

	fmt.Printf("Created access token ID: %s\n", id.String())
	fmt.Printf("Token: %s\n", token)

	return nil
}

var tokenDeleteCmd = &cobra.Command{
	Use:   "delete [token ID]",
	Short: "Delete a token",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenDeleteCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenDeleteCmd(cctx *TokenCommandContext) error {

	tokenID := cctx.Args[0]

	err := cctx.Token.GetClient().DeleteToken(tokenID)
	if err != nil {
		return err
	}

	fmt.Printf("Token \"%s\" was deleted\n", tokenID)

	return nil
}

var tokenUpdateCmd = &cobra.Command{
	Use:   "update [token ID]",
	Short: "Update a token",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenUpdateCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenUpdateCmd(cctx *TokenCommandContext) error {

	tokenID := cctx.Args[0]
	changed := false

	// Getting token information
	token, err := cctx.Token.GetClient().GetToken(tokenID)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found token \"%s\"\n", tokenID))
	}

	// Update description
	if cctx.Cmd.Flags().Changed("desc") {
		token.Description = tokenDesc
		changed = true
	}

	// Update enabled
	if cctx.Cmd.Flags().Changed("enabled") {
		token.Enabled = tokenEnabled
		changed = true
	}

	// Nothing's changed
	if !changed {
		return nil
	}

	// Update
	_, err = cctx.Token.GetClient().UpdateToken(tokenID, token)
	if err != nil {
		return err
	}

	return nil
}

var tokenInfoCmd = &cobra.Command{
	Use:   "info [token ID]",
	Short: "Show information about token",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenInfoCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenInfoCmd(cctx *TokenCommandContext) error {

	tokenID := cctx.Args[0]

	// Getting token information
	token, err := cctx.Token.GetClient().GetToken(tokenID)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found token \"%s\"\n", tokenID))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	//	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	//table.SetNoWhiteSpace(true)

	// Basic information
	table.Append([]string{
		"ID:",
		token.ID,
	})

	table.Append([]string{
		"Description:",
		token.Description,
	})

	var status string
	if token.Enabled {
		status = "enabled"
	} else {
		status = "disabled"
	}

	table.Append([]string{
		"Status:",
		status,
	})

	table.Append([]string{
		"Updated:",
		token.UpdatedAt.String(),
	})

	table.Append([]string{
		"Created:",
		token.CreatedAt.String(),
	})

	fmt.Printf("Information for Token %s\n\n", tokenID)
	fmt.Printf("Configuration:\n\n")

	table.Render()

	fmt.Println("")

	// Permissions
	if len(token.Permissions) > 0 {
		for permission, _ := range token.Permissions {
			fmt.Println("")
			fmt.Println("Permissions:")
			fmt.Println("")

			fmt.Println(permission)
		}

		fmt.Println("")
	}

	return nil
}

var tokenGrantCmd = &cobra.Command{
	Use:   "grant [token ID] [permission]",
	Short: "Grant permission to specific token",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenGrantCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenGrantCmd(cctx *TokenCommandContext) error {

	tokenID := cctx.Args[0]
	permission := cctx.Args[1]

	// Getting token information
	token, err := cctx.Token.GetClient().GetToken(tokenID)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found token \"%s\"\n", tokenID))
	}

	// Permission exists already
	if token.Permissions == nil {
		token.Permissions = make(map[string]*token_sdk.Permission)
	}

	if _, ok := token.Permissions[permission]; ok {
		return nil
	}

	// Add permission to token setting
	token.Permissions[permission] = &token_sdk.Permission{}

	// Update
	_, err = cctx.Token.GetClient().UpdateToken(tokenID, token)
	if err != nil {
		return err
	}

	return nil
}

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke [token ID] [permission]",
	Short: "Revoke permission from specific token",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runTokenCmd(runTokenRevokeCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runTokenRevokeCmd(cctx *TokenCommandContext) error {

	tokenID := cctx.Args[0]
	permission := cctx.Args[1]

	// Getting token information
	token, err := cctx.Token.GetClient().GetToken(tokenID)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return errors.New(fmt.Sprintf("Not found token \"%s\"\n", tokenID))
	}

	// Permission exists already
	if token.Permissions == nil {
		return nil
	}

	if _, ok := token.Permissions[permission]; !ok {
		return nil
	}

	// Revoke permission
	delete(token.Permissions, permission)

	// Update
	_, err = cctx.Token.GetClient().UpdateToken(tokenID, token)
	if err != nil {
		return err
	}

	return nil
}
