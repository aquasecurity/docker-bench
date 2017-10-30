// Copyright © 2017 Aqua Security Software Ltd. <info@aquasec.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	goflag "flag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dbDir = "appdb"

	cfgFile   string
	checkList string
	name      string
	jsonFmt   bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "docker-bench",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: app,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	goflag.Set("logtostderr", "true")
	goflag.CommandLine.Parse([]string{})

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docker-bench.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringVarP(&name, "name", "n", "", "name of app to run check against")
	viper.BindPFlag("name", RootCmd.Flags().Lookup("name"))

	RootCmd.Flags().StringVarP(&dbDir, "dbdir", "d", "appdb", "directory to get benchmark definitions")
	viper.BindPFlag("dbdir", RootCmd.Flags().Lookup("dbdir"))

	RootCmd.PersistentFlags().BoolVar(&jsonFmt, "json", false, "Prints the results as JSON")
	viper.BindPFlag("json", RootCmd.Flags().Lookup("json"))

	RootCmd.PersistentFlags().StringVarP(
		&checkList,
		"check",
		"c",
		"",
		`A comma-delimited list of checks to run as specified in CIS document. Example --check="1.1.1,1.1.2"`,
	)

	goflag.CommandLine.VisitAll(func(goflag *goflag.Flag) {
		RootCmd.PersistentFlags().AddGoFlag(goflag)
	})

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".docker-bench") // name of config file (without extension)
	viper.AddConfigPath("$HOME")         // adding home directory as first search path
	viper.AutomaticEnv()                 // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
