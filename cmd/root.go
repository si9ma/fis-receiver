/*
Copyright © 2020 si9ma <si9ma@si9ma.com>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "fis-receiver",
	Short: "a fis receiver with golang",
	Long:  `a fis receiver build with golang`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

