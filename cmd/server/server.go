/*
Copyright 2020 Christian Banse

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

package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/oxisto/titan"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/datafetch"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/routes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	RedisFlag              = "redis"
	PostgresFlag           = "postgres"
	ListenFlag             = "listen"
	CorporationIDFlag      = "corporationID"
	CacheManufacturingFlag = "cache.manufacturing"
	EveClientID            = "eve.clientID"
	EveSecretKey           = "eve.secretKey"
	EveRedirectURI         = "eve.redirectURI"

	DefaultRedis              = "localhost:6379"
	DefaultPostgres           = "localhost"
	DefaultListen             = ":4300"
	DefaultCorporationID      = 0
	DefaultCacheManufacturing = "true"
	DefaultEmpty              = ""

	EnvPrefix = "TITAN"
)

var serverCmd = &cobra.Command{
	Use:   "titan-server",
	Short: "titan-server is the main API server for Titan",
	Long:  "Titan Server is the main component of Titan. It takes care of computing all of your manufacturing needs",
	Run:   doCmd,
}

func init() {
	cobra.OnInitialize(initConfig)

	serverCmd.Flags().String(ListenFlag, DefaultListen, "Host and port to listen to")
	serverCmd.Flags().String(RedisFlag, DefaultRedis, "Host and port of redis server")
	serverCmd.Flags().String(PostgresFlag, DefaultPostgres, "Connection string for PostgreSQL")
	serverCmd.Flags().Int32(CorporationIDFlag, DefaultCorporationID, "If specified, limits access to this corporation ID")
	serverCmd.Flags().String(EveClientID, DefaultEmpty, "The EVE SSO Client ID")
	serverCmd.Flags().String(EveSecretKey, DefaultEmpty, "The EVE SSO Secret Key")
	serverCmd.Flags().String(EveRedirectURI, DefaultEmpty, "The EVE SSO Redirect URI")

	// TODO: this should actually be a bool but they behave wierdly
	serverCmd.Flags().String(CacheManufacturingFlag, DefaultCacheManufacturing, "Specifies, whether to regularly cache manufacturing during the runtime of the server")

	viper.BindPFlag(ListenFlag, serverCmd.Flags().Lookup(ListenFlag))
	viper.BindPFlag(RedisFlag, serverCmd.Flags().Lookup(RedisFlag))
	viper.BindPFlag(PostgresFlag, serverCmd.Flags().Lookup(PostgresFlag))
	viper.BindPFlag(CorporationIDFlag, serverCmd.Flags().Lookup(CorporationIDFlag))
	viper.BindPFlag(CacheManufacturingFlag, serverCmd.Flags().Lookup(CacheManufacturingFlag))
	viper.BindPFlag(EveClientID, serverCmd.Flags().Lookup(EveClientID))
	viper.BindPFlag(EveSecretKey, serverCmd.Flags().Lookup(EveSecretKey))
	viper.BindPFlag(EveRedirectURI, serverCmd.Flags().Lookup(EveRedirectURI))
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// TODO: should we read config here ?!
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Could not read config: %s", err)
	}
}

func doCmd(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
	log.Info("Starting server...")

	if !cache.InitSSO(
		viper.GetString(EveClientID),
		viper.GetString(EveSecretKey),
		viper.GetString(EveRedirectURI)) {
		log.Errorf("Could not initialize EVE SSO, please specify a client ID, secret and redirect URI.")
		return
	}

	if err := cache.InitCache(viper.GetString(RedisFlag)); err != nil {
		log.Errorf("Could not initialize cache: %s", err)
		return
	}

	db.InitPostgreSQL(viper.GetString(PostgresFlag))

	app := titan.App{
		CacheManufacturing: viper.GetBool(CacheManufacturingFlag),
		CorporationID:      int32(viper.GetInt(CorporationIDFlag)),
	}

	app.ImportSDE()

	go app.ServerLoop()
	go app.JournalLoop()

	var division int32

	for division = 1; division <= 3; division++ {
		transactionFetcher := datafetch.NewTransactionFetcher(app.CorporationID, division)
		go transactionFetcher.StartLoop()
	}

	walletFetcher := datafetch.NewWallerFetcher(app.CorporationID, division)
	go walletFetcher.StartLoop()

	//go app.TransactionLoop()
	//go ContractsLoop()

	router := routes.NewRouter(int32(viper.GetInt(CorporationIDFlag)))
	err := http.ListenAndServe(viper.GetString(ListenFlag), router)

	log.Errorf("An error occured: %v", err)
}

func main() {
	log.SetLevel(log.DebugLevel)

	if err := serverCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
