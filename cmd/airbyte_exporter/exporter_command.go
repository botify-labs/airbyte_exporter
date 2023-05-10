// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/virtualtam/venom"
)

const (
	defaultDatabaseAddr     string = "localhost:5432"
	defaultDatabaseName     string = "airbyte"
	defaultDatabaseUser     string = "airbyte_exporter"
	defaultDatabasePassword string = "airbyte_exporter"
)

var (
	defaultLogLevelValue string = zerolog.LevelInfoValue
	logLevelValue        string

	logLevelValues = []string{
		zerolog.LevelTraceValue,
		zerolog.LevelDebugValue,
		zerolog.LevelInfoValue,
		zerolog.LevelWarnValue,
		zerolog.LevelErrorValue,
		zerolog.LevelFatalValue,
		zerolog.LevelPanicValue,
	}

	databaseAddr     string
	databaseName     string
	databaseUser     string
	databasePassword string
)

// NewExporterCommand initializes the exporter's CLI entrypoint and command flags.
func NewExporterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "airbyte_exporter",
		Short: "Airbyte Exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Configuration file lookup paths
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			homeConfigPath := filepath.Join(home, ".config")

			configPaths := []string{DefaultConfigPath, homeConfigPath, "."}

			// Inject global configuration as a pre-run hook
			//
			// This is required to let Viper load environment variables and
			// configuration entries before invoking nested commands.
			if err := venom.Inject(cmd, EnvPrefix, configPaths, ConfigName, false); err != nil {
				return err
			}

			// Global logger configuration
			var logLevel zerolog.Level

			if err := logLevel.UnmarshalText([]byte(logLevelValue)); err != nil {
				log.Error().Err(err).Msg("invalid log level")
				return err
			}

			log.Info().Str("log_level", logLevelValue).Msg("setting log level")
			zerolog.SetGlobalLevel(logLevel)

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(
		&logLevelValue,
		"log-level",
		defaultLogLevelValue,
		fmt.Sprintf(
			"Log level (%s)",
			strings.Join(logLevelValues, ", "),
		),
	)

	cmd.PersistentFlags().StringVar(
		&databaseAddr,
		"db-addr",
		defaultDatabaseAddr,
		"Database address (host:port)",
	)
	cmd.PersistentFlags().StringVar(
		&databaseName,
		"db-name",
		defaultDatabaseName,
		"Database name",
	)
	cmd.PersistentFlags().StringVar(
		&databaseUser,
		"db-user",
		defaultDatabaseUser,
		"Database user",
	)
	cmd.PersistentFlags().StringVar(
		&databasePassword,
		"db-password",
		defaultDatabasePassword,
		"Database password",
	)

	return cmd
}
