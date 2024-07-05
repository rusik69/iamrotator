package main

import (
	"os"

	"github.com/rusik69/iamrotator/pkg/args"
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logLevelString := os.Getenv("LOG_LEVEL")
	if logLevelString == "" {
		logLevelString = "info"
	}
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err != nil {
		logrus.Warnf("Invalid log level '%s', defaulting to 'info'", logLevelString)
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	arg := args.Parse()
	cfg, err := config.Load(arg.ConfigPath)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info("Loaded configuration from", arg.ConfigPath)
	logrus.Debugf("Config: %v\n", cfg)

	awsSess, err := aws.CreateSession(cfg.AWS)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info("Checking AWS account", cfg.AWS.AccountID)
	if arg.Action == "createuser" {
		newKeyID, NewKeySecret, err := aws.CheckOrCreateIamUser(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Access Key ID:", newKeyID)
		logrus.Info("Secret Access Key:", NewKeySecret)
	} else if arg.Action == "createstackset" {
		logrus.Info("Checking or creating stack set")
		err = aws.CheckOrCreateStackSet(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
	} else if arg.Action == "removeuser" {
		logrus.Info("Removing IAM user")
		err = aws.RemoveIamUser(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("IAM user removed")
	} else if arg.Action == "removestackset" {
		logrus.Info("Emptying stack set")
		err = aws.EmptyStackSet(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Removing stack set")
		err = aws.RemoveStackSet(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Stack set removed")
	} else {
		logrus.Info("Usage: iamrotator <action> <configpath>")
		os.Exit(1)
	}
}
