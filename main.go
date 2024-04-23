/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/avinashsivaraman/gq/cmd"
	"github.com/spf13/viper"
)

func init() {
	// Set the config file name and file path
	viper.SetConfigName(".gq")   // Specify the config file name without extension
	viper.SetConfigType("yaml")  // Set the type of config file
	viper.AddConfigPath("$HOME/.config/gq") // Add the root directory to search for the config file

	// If needed, you can also set default config options here

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		// Handle errors such as config file not found
		// You can choose to create a default config file here if needed
		fmt.Printf("Error reading config file: %s \n", err)
		os.Exit(1)
	}
}

func main() {
	cmd.Execute()
}
