package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mpppk/connect-to-gce-win/lib"
	"github.com/mpppk/hlb/hlblib"

	"github.com/atotto/clipboard"
	"github.com/mpppk/gce-auto-connect/password"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"

	daisyCompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "connect-to-gce-win",
	Short: "connect to windows on GCE via RDP",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var config lib.Config
		err := viper.Unmarshal(&config)
		if err != nil {
			panic(errors.Wrap(err, "failed to unmarshal config"))
		}

		configDirPath, err := lib.GetConfigDirPath()
		if err != nil {
			panic(errors.Wrap(err, "failed to get config dir path"))
		}
		rdpFilePath := filepath.Join(configDirPath, "connect-to-gce-win.rdp")

		ctx := context.Background()
		service, err := daisyCompute.NewClient(ctx)

		if err != nil {
			panic(err)
		}
		instance, err := service.GetInstance(config.Project, config.Zone, config.InstanceName)

		if err != nil {
			panic(errors.Wrap(err, "failed to get instance"))
		}

		if instance.Status == "TERMINATED" {
			fmt.Print("Starting instance")
			err = service.StartInstance(config.Project, config.Zone, config.InstanceName)
			if err != nil {
				panic(errors.Wrap(err, "failed to start instance"))
			}
		}

		instance, err = service.GetInstance(config.Project, config.Zone, config.InstanceName)

		if err != nil {
			panic(errors.Wrap(err, "failed to get instance"))
		}

		natIP, err := lib.ExtractNatIpFromInstance(instance)
		if err != nil {
			panic(errors.Wrap(err, "failed to get NAT IP"))
		}

		err = lib.GenerateRDPFile(rdpFilePath, natIP, "niboshiporipori")
		if err != nil {
			panic(errors.Wrap(err, "failed to generate RDP file"))
		}

		newPassword, err := password.ResetPassword(service, config.InstanceName, config.Zone, config.Project, config.UserName)
		err = clipboard.WriteAll(newPassword)
		if err != nil {
			panic(errors.Wrap(err, "failed to copy password to clipboard"))
		}
		fmt.Println("new password has been copied to clipboard")
		err = open.Run(rdpFilePath)
		if err != nil {
			panic(errors.Wrap(err, "failed to open RDP file"))
		}
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/connect-to-gce-win/.connect-to-gce-win.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".connect-to-gce-win")
		configFilePath, err := hlblib.GetConfigDirPath()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(configFilePath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
