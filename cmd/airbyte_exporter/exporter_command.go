// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/virtualtam/venom"

	"github.com/botify-labs/airbyte_exporter/v2/internal/airbyte"
)

const (
	defaultListenAddr string = "0.0.0.0:8080"

	defaultDatabaseAddr     string = "localhost:5432"
	defaultDatabaseSSLMode  string = "disable"
	defaultDatabaseName     string = "airbyte"
	defaultDatabaseUser     string = "airbyte_exporter"
	defaultDatabasePassword string = "airbyte_exporter"

	databaseDriver string = "pgx"
)

var (
	listenAddr           string
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
	databaseSSLMode  string
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

			// Encode the database password with percent encoding in case it contains special characters.
			// https://www.postgresql.org/docs/current/libpq-connect.html
			// https://datatracker.ietf.org/doc/html/rfc3986#section-2.1
			databasePassword = url.QueryEscape(databasePassword)
			databaseURI := fmt.Sprintf(
				"postgres://%s:%s@%s/%s?sslmode=%s",
				databaseUser,
				databasePassword,
				databaseAddr,
				databaseName,
				databaseSSLMode,
			)

			// Database connection pool
			pgxPool, err := pgxpool.New(context.Background(), databaseURI)
			if err != nil {
				log.Error().
					Err(err).
					Str("database_driver", databaseDriver).
					Str("database_addr", databaseAddr).
					Str("database_name", databaseName).
					Msg("database: failed to create connection pool")
				return err
			}

			if err := pgxPool.Ping(context.Background()); err != nil {
				log.Error().
					Err(err).
					Str("database_driver", databaseDriver).
					Str("database_addr", databaseAddr).
					Str("database_name", databaseName).
					Msg("database: failed to ping")
				return err
			}

			log.Info().
				Str("database_driver", databaseDriver).
				Str("database_addr", databaseAddr).
				Str("database_name", databaseName).
				Msg("database: successfully created connection pool")

			// Airbyte Exporter services
			airbyteRepository := airbyte.NewRepository(pgxPool)
			airbyteService := airbyte.NewService(airbyteRepository)

			httpServer := newServer(airbyteService, listenAddr)

			log.Info().Str("addr", listenAddr).Msg("starting HTTP server")
			return httpServer.ListenAndServe()
		},
	}

	cmd.Flags().StringVar(
		&listenAddr,
		"listen-addr",
		defaultListenAddr,
		"Listen to this address (host:port)",
	)

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
		&databaseSSLMode,
		"db-sslmode",
		defaultDatabaseSSLMode,
		"Database sslmode",
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
