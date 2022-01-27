/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// passwordCmd represents the password command.
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter password (it will be hidden): ")
		password, err := term.ReadPassword(0)
		if err != nil {
			log.Fatal("error reading password")
		}

		passwordHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("error hashing password")
		}

		fmt.Printf("%s\n", passwordHash)
	},
}

func init() {
	rootCmd.AddCommand(passwordCmd)
}
