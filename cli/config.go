package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup config file for rattlesnakeos-stack",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan(fmt.Sprintln("Stack name is used as an identifier for all the AWS components that get deployed. This name must be unique or stack deployment will fail."))
		validate := func(input string) error {
			if len(input) < 1 {
				return errors.New("Stack name is too short")
			}
			return nil
		}
		namePrompt := promptui.Prompt{
			Label:    "Stack name ",
			Validate: validate,
			Default:  viper.GetString("name"),
		}
		result, err := namePrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		viper.Set("name", result)

		color.Cyan(fmt.Sprintf("Stack region is the AWS region where you would like to deploy your stack. Valid options: %v\n",
			strings.Join(supportedRegions, ", ")))
		validate = func(input string) error {
			if len(input) < 1 {
				return errors.New("Stack region is too short")
			}
			found := false
			for _, region := range supportedRegions {
				if input == region {
					found = true
					break
				}
			}
			if !found {
				return errors.New("Invalid region")
			}
			return nil
		}
		regionPrompt := promptui.Prompt{
			Label:    "Stack region ",
			Default:  viper.GetString("region"),
			Validate: validate,
		}
		result, err = regionPrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		viper.Set("region", result)

		color.Cyan(fmt.Sprintln("Device is the device codename (e.g. sailfish). Supported devices:", supportDevicesOutput))
		validate = func(input string) error {
			if len(input) < 1 {
				return errors.New("Device name is too short")
			}
			found := false
			for _, d := range supportedDevicesCodename {
				if input == d {
					found = true
					break
				}
			}
			if !found {
				return errors.New("Invalid device")
			}
			return nil
		}
		devicePrompt := promptui.Prompt{
			Label:    "Device ",
			Default:  viper.GetString("device"),
			Validate: validate,
		}
		result, err = devicePrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		viper.Set("device", result)

		color.Cyan(fmt.Sprintln("Email address you would like to send build notifications to."))
		validate = func(input string) error {
			if !strings.Contains(input, "@") {
				return errors.New("Must provide valid email")
			}
			return nil
		}
		emailPrompt := promptui.Prompt{
			Label:    "Email ",
			Validate: validate,
			Default:  viper.GetString("email"),
		}
		result, err = emailPrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		viper.Set("email", result)

		color.Cyan(fmt.Sprintln("SSH keypair name is the name of your EC2 keypair that was generated/uploaded in AWS."))
		validate = func(input string) error {
			if len(input) < 1 {
				return errors.New("SSH keypair name is too short")
			}
			return nil
		}
		keypairPrompt := promptui.Prompt{
			Label:    "SSH Keypair Name ",
			Default:  viper.GetString("ssh-key"),
			Validate: validate,
		}
		result, err = keypairPrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		viper.Set("ssh-key", result)

		err = viper.WriteConfigAs(configFileFullPath)
		if err != nil {
			log.WithError(err).Fatalf("failed to write config file %s", configFileFullPath)
		}
		log.Infof("rattlesnakeos-stack config file has been written to %v", configFileFullPath)
	},
}
