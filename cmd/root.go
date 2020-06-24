package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tionebsalocin/vault-wrapper/runner"
)

var path []string
var vault_addr string
var auth_path string
var user string
var realm string
var spn string
var krb5 string
var keytab string
var prefix string
var token_only bool
var export_only bool

func init() {
	RootCmd.PersistentFlags().StringSliceVarP(&path, "secret_path", "p", []string{}, "Secret path in Vault path/to/secret:key")
	RootCmd.PersistentFlags().StringVarP(&vault_addr, "vault_addr", "v", "", "Vault address")
	RootCmd.PersistentFlags().StringVarP(&auth_path, "auth_path", "a", "", "login url for vault kerberos plugin")
	RootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user login")
	RootCmd.PersistentFlags().StringVarP(&realm, "realm", "r", "", "kerberos realm for user")
	RootCmd.PersistentFlags().StringVarP(&spn, "spn", "s", "", "Vault Service Principal Name")
	RootCmd.PersistentFlags().StringVarP(&krb5, "krb5", "k", "/etc/krb5.conf", "Path to krb5.conf")
	RootCmd.PersistentFlags().StringVarP(&keytab, "keytab", "t", "", "Path to keytab file")
	RootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "e", "SECRET_", "Environment variable prefix")
	RootCmd.PersistentFlags().BoolVarP(&token_only, "token-only", "o", false, "Only output vault token (no execution wrapped)")
	RootCmd.PersistentFlags().BoolVarP(&export_only, "export-only", "x", false, "Only output command to export env variables (no execution wrapped)")
}

var RootCmd = &cobra.Command{
	Use:  "vault-wrapper",
	Long: "Collect secret from vault",
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var exec_command []string
		if len(args) > 0 {
			exec_command = []string{"sh", "-c", args[0]}
		} else if token_only == false && export_only == false {
			log.Fatal("[VaultWrapper] Syntax Error: Missing command line to wrap")
		}
		runner.Run(runner.Config{
			VaultPath:      path,
			VaultAddr:      vault_addr,
			VaultAuthPath:  auth_path,
			VaultUser:      user,
			VaultUserRealm: realm,
			VaultSPN:       spn,
			Krb5:           krb5,
			Keytab:         keytab,
			EnvPrefix:      prefix,
			TokenOnly:      token_only,
			ExportOnly:     export_only,
			CommandLine:    exec_command,
		})
	},
}
