package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
	"github.com/k0kubun/pp"
	"github.com/mpppk/gce-auto-connect/password"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"

	daisyCompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/compute/v1"
)

var cfgFile string
var project string
var zone string
var name string

var userName = "niboshiporipori"
var rdpFilePath = "./gce-windows.rdp"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "connect-to-gce-win",
	Short: "connect to windows on GCE via RDP",
	//Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service, err := daisyCompute.NewClient(ctx)

		if err != nil {
			panic(err)
		}
		instance, err := service.GetInstance(project, zone, name)

		if err != nil {
			panic(errors.Wrap(err, "failed to get instance"))
		}

		if instance.Status == "TERMINATED" {
			fmt.Print("starting instance...")
			err = service.StartInstance(project, zone, name)
			if err != nil {
				panic(errors.Wrap(err, "failed to start instance"))
			}
		}

		instance, err = service.GetInstance(project, zone, name)

		if err != nil {
			panic(errors.Wrap(err, "failed to get instance"))
		}

		pp.Println(instance)

		natIP, err := extractNatIpFromInstance(instance)
		if err != nil {
			panic(errors.Wrap(err, "failed to get NAT IP"))
		}

		err = generateRDPFile(rdpFilePath, natIP, "niboshiporipori")
		if err != nil {
			panic(errors.Wrap(err, "failed to generate RDP file"))
		}

		newPassword, err := password.ResetPassword(service, name, zone, project, userName)
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.connect-to-gce-win.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVar(&project, "project", "windows-provisioning", "Project name on GCP")
	rootCmd.Flags().StringVar(&zone, "zone", "asia-northeast1-b", "Zone of GCP")
	rootCmd.Flags().StringVar(&name, "name", "magic-arena", "Instance name to create")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		// Search config in home directory with name ".connect-to-gce-win" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".connect-to-gce-win")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func generateRDPFile(filePath, ip, userName string) error {
	contents := "full address:s:" + ip + "\n" + "username:s:" + userName
	return ioutil.WriteFile(filePath, []byte(contents), 0644)
}

func extractNatIpFromInstance(instance *compute.Instance) (string, error) {
	if len(instance.NetworkInterfaces) == 0 {
		return "", errors.New("no NetworkInterfaces found")
	}
	accessConfigs := instance.NetworkInterfaces[0].AccessConfigs
	if len(accessConfigs) == 0 {
		return "", errors.New("no AccessConfigs found")
	}
	return accessConfigs[0].NatIP, nil
}
