package cmd

import (
	"errors"
	"os"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
	"github.com/ajgon/mailbowl/listener"
	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/ajgon/mailbowl/process"
	"github.com/ajgon/mailbowl/relay"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "mailbowl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		tlsConfig, err := smtp.NewTLS(
			viper.GetString("smtp.tls.key"),
			viper.GetString("smtp.tls.certificate"),
			viper.GetString("smtp.tls.key_file"),
			viper.GetString("smtp.tls.certificate_file"),
			viper.GetBool("smtp.tls.force_for_starttls"),
		)
		if err != nil && !errors.Is(err, smtp.ErrTLSNotConfigured) {
			log.Fatalf("invalid TLS config: %s", err.Error())
		}

		outgoingServer, err := relay.NewOutgoingServer(
			viper.GetString("relay.outgoing_server.address"),
			viper.GetString("relay.outgoing_server.auth_method"),
			viper.GetString("relay.outgoing_server.connection_type"),
			viper.GetString("relay.outgoing_server.from_email"),
			viper.GetString("relay.outgoing_server.password"),
			viper.GetString("relay.outgoing_server.username"),
			viper.GetBool("relay.outgoing_server.verify_tls"),
		)
		if err != nil {
			log.Fatalf("invalid outgoing server config: %s", err.Error())
		}

		smtpAuthUsers := make([]*smtp.AuthUser, 0)

		if users, ok := viper.Get("smtp.auth.users").([]map[string]string); ok {
			for _, user := range users {
				if user["email"] != "" && user["password_hash"] != "" {
					smtpAuthUsers = append(smtpAuthUsers, smtp.NewAuthUser(user["email"], user["password_hash"]))
				}
			}
		}

		httpServer := listener.NewHTTP()
		smtpServer := smtp.NewSMTP(
			smtp.NewAuth(viper.GetBool("smtp.auth.enabled"), smtpAuthUsers),
			viper.GetString("smtp.hostname"),
			smtp.NewLimit(
				viper.GetInt("smtp.limit.connections"),
				viper.GetInt("smtp.limit.message_size"),
				viper.GetInt("smtp.limit.recipients"),
			),
			smtp.NewTimeout(
				viper.GetString("smtp.timeout.read"),
				viper.GetString("smtp.timeout.write"),
				viper.GetString("smtp.timeout.data"),
			),
			tlsConfig,
			smtp.NewWhitelist(
				viper.GetStringSlice("smtp.whitelist.cidrs"),
			),
			viper.GetStringSlice("smtp.listen"),
			relay.NewRelay(outgoingServer),
		)

		manager := process.NewManager()
		manager.AddListener(httpServer)
		manager.AddListener(smtpServer)

		manager.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.CobraInitialize(cfgFile))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./mailbowl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
