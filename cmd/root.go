/*
Copyright © 2022 kockicica@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"lpfr-o-matic/pkg/sys"
	"lpfr-o-matic/pkg/watchdog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	exePath       string
	checkUrl      string
	interval      int
	pin           string
	noAutoPin     bool
	middlewareApp string
)

var rootCmd = &cobra.Command{
	Use:     "lpfr-o-matic",
	Short:   "LPFR runner & monitor",
	Version: "0.1.12",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := sys.CreateMutex("lpfr-o-matic")
		if err != nil {
			fmt.Println("It looks like another instance of lpfr-o-matic is running")
			return err
		}

		wdConfig := watchdog.WatchdogConfig{}
		if err := viper.Unmarshal(&wdConfig); err != nil {
			return err
		}
		wd := watchdog.NewWatchdog(wdConfig)
		err = wd.Start()
		if err != nil {
			log.Fatalln(err)
			return err
		}
		wait := make(chan os.Signal)
		log.Println("Started")
		signal.Notify(wait, os.Kill, os.Interrupt)
		<-wait
		log.Println("Stopped")
		return nil
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.lpfr-o-matic.yaml)")
	rootCmd.Flags().StringVar(&exePath, "exepath", "lpfr.lnk", "LPFR executable name / shortcut")
	rootCmd.Flags().StringVar(&checkUrl, "checkurl", "http://localhost:7555", "url of the LPFR to check")
	rootCmd.Flags().IntVar(&interval, "interval", 10, "interval (in seconds) to perform checks")
	rootCmd.Flags().StringVar(&pin, "pin", "", "smart card pin code")
	rootCmd.Flags().BoolVar(&noAutoPin, "nopin", false, "skip automatic pin setup")
	rootCmd.Flags().StringVar(&middlewareApp, "middleware", "", "middleware application to start on successful status")
	rootCmd.Flags().Bool("telegram", false, "send telegram status messages")
	rootCmd.Flags().String("telegram-api-key", "", "telegram bot api key")
	rootCmd.Flags().String("telegram-chat-id", "", "telegram chat id")
	rootCmd.Flags().String("telegram-sender", "", "sender identification")

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatalln(err)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in working directory with name ".lpfr-o-matic" (without extension).
		viper.AddConfigPath(wd)
		viper.SetConfigName(".lpfr-o-matic")
	}

	viper.SetEnvPrefix("LPFR_")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
