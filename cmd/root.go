package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
	"github.com/ajgon/mailbowl/listener"
	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/ajgon/mailbowl/process"
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

		smtpAuthUsers := make([]*smtp.AuthUser, 0)

		if users, ok := viper.Get("smtp.auth.users").([]map[string]string); ok {
			for _, user := range users {
				smtpAuthUsers = append(smtpAuthUsers, smtp.NewAuthUser(user["email"], user["password_hash"]))
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
			viper.GetStringSlice("smtp.listen"),
			smtp.NewTimeout(
				viper.GetString("smtp.timeout.read"),
				viper.GetString("smtp.timeout.write"),
				viper.GetString("smtp.timeout.data"),
			),
			tlsConfig,
			smtp.NewWhitelist(
				viper.GetStringSlice("smtp.whitelist.cidrs"),
			),
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mailbowl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.SetDefaults()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mailbowl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mailbowl")
	}

	viper.SetEnvPrefix("MAILBOWL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	err := config.ConfigureLogger(os.Stdout, os.Stderr)
	cobra.CheckErr(err)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	viperSettingsJSON, err := json.Marshal(viper.AllSettings())
	cobra.CheckErr(err)
	log.Debug("Loaded config: ", string(viperSettingsJSON))
}
