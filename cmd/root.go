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
	"google.golang.org/api/compute/v1"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "connect-to-gce-win",
	Short: "connect to windows on GCE via RDP",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var config lib.Config
		err := viper.Unmarshal(&config)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to unmarshal config"))

		configDirPath, err := lib.GetConfigDirPath()
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to get config dir path"))

		rdpFilePath := filepath.Join(configDirPath, "connect-to-gce-win.rdp")

		ctx := context.Background()
		service, err := daisyCompute.NewClient(ctx)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to create client of GCE"))

		var instance *compute.Instance
		if config.InstanceName != "" {
			instance, err = service.GetInstance(config.Project, config.Zone, config.InstanceName)
			lib.PanicIfErrorExist(errors.Wrap(err, "failed to get instance"))
		} else {
			instance, err = lib.ChooseInstance(config.Project, config.Zone)
			lib.PanicIfErrorExist(errors.Wrap(err, "failed to get instance"))
		}

		fmt.Println("Trying to connect to " + instance.Name)

		if instance.Status == "TERMINATED" {
			fmt.Println("Starting instance")
			err = service.StartInstance(config.Project, config.Zone, instance.Name)
			lib.PanicIfErrorExist(errors.Wrap(err, "failed to start instance"))
		}

		instance, err = service.GetInstance(config.Project, config.Zone, instance.Name)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to get instance"))

		natIP, err := lib.ExtractNatIpFromInstance(instance)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to get NAT IP"))

		err = lib.GenerateRDPFile(rdpFilePath, natIP, config.UserName)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to generate RDP file"))

		newPassword, err := password.ResetPassword(service, instance.Name, config.Zone, config.Project, config.UserName)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to reset password"))

		err = clipboard.WriteAll(newPassword)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to copy password to clipboard"))
		fmt.Println("New password has been copied to clipboard")

		err = open.Run(rdpFilePath)
		lib.PanicIfErrorExist(errors.Wrap(err, "failed to open RDP file"))
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
