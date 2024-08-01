package main

import (
	"github.com/rusik69/iamrotator/pkg/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		err := web.Listen(configPath)
		if err != nil {
			logrus.Panic(err)
		}
	},
}
