// Copyright 2020 Kentaro Hibino. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"github.com/cloudjjcc/asynq"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/cloudjjcc/asynq/internal/base"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// Flags
var uri string
var db int
var password string
var prefix string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "asynq",
	Short:   "A monitoring tool for asynq queues",
	Long:    `Asynq is a montoring CLI to inspect tasks and queues managed by asynq.`,
	Version: base.Version,
}

var versionOutput = fmt.Sprintf("asynq version %s\n", base.Version)

var versionCmd = &cobra.Command{
	Use:    "version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(versionOutput)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)
	rootCmd.SetVersionTemplate(versionOutput)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file to set flag defaut values (default is $HOME/.asynq.yaml)")
	rootCmd.PersistentFlags().StringVarP(&uri, "uri", "u", "127.0.0.1:6379", "redis server URI")
	rootCmd.PersistentFlags().IntVarP(&db, "db", "n", 0, "redis database number (default is 0)")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password to use when connecting to redis server")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "test:", "redis prefix")
	viper.BindPFlag("uri", rootCmd.PersistentFlags().Lookup("uri"))
	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if prefix != "" {
		asynq.SetRedisBasePrefix(prefix)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".asynq" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".asynq")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// printTable is a helper function to print data in table format.
//
// cols is a list of headers and printRow specifies how to print rows.
//
// Example:
// type User struct {
//     Name string
//     Addr string
//     Age  int
// }
// data := []*User{{"user1", "addr1", 24}, {"user2", "addr2", 42}, ...}
// cols := []string{"Name", "Addr", "Age"}
// printRows := func(w io.Writer, tmpl string) {
//     for _, u := range data {
//         fmt.Fprintf(w, tmpl, u.Name, u.Addr, u.Age)
//     }
// }
// printTable(cols, printRows)
func printTable(cols []string, printRows func(w io.Writer, tmpl string)) {
	format := strings.Repeat("%v\t", len(cols)) + "\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	var headers []interface{}
	var seps []interface{}
	for _, name := range cols {
		headers = append(headers, name)
		seps = append(seps, strings.Repeat("-", len(name)))
	}
	fmt.Fprintf(tw, format, headers...)
	fmt.Fprintf(tw, format, seps...)
	printRows(tw, format)
	tw.Flush()
}
