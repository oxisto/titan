/*
Copyright 2018 Christian Banse

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
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/oxisto/titan/manufacturing"

	"fmt"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
	"github.com/oxisto/titan/routes"
	"github.com/oxisto/titan/slack"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	RedisFlag              = "redis"
	MongoFlag              = "mongo"
	ListenFlag             = "listen"
	CorporationIDFlag      = "corporationID"
	SlackAPITokenFlag      = "slack.token"
	CacheManufacturingFlag = "cache.manufacturing"
	EveClientID            = "eve.clientID"
	EveSecretKey           = "eve.secretKey"
	EveRedirectURI         = "eve.redirectURI"

	DefaultRedis              = "localhost:6379"
	DefaultMongo              = "localhost:27017"
	DefaultListen             = ":4300"
	DefaultCorporationID      = 0
	DefaultSlackAPIToken      = DefaultEmpty
	DefaultCacheManufacturing = "true"
	DefaultEmpty              = ""

	EnvPrefix = "TITAN"
)

// DebugLogWriter implements io.Writer and writes all incoming text out to log level info.
type DebugLogWriter struct {
	Component string
}

func (d DebugLogWriter) Write(p []byte) (n int, err error) {
	log.WithField("component", d.Component).Debug(strings.TrimRight(string(p), "\n"))

	return len(p), nil
}

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
	serverCmd.Flags().String(MongoFlag, DefaultMongo, "Host and port of mongodb server")
	serverCmd.Flags().Int32(CorporationIDFlag, DefaultCorporationID, "If specified, limits access to this corporation ID")
	serverCmd.Flags().String(SlackAPITokenFlag, DefaultSlackAPIToken, "The token for Slack integration")
	serverCmd.Flags().String(EveClientID, DefaultEmpty, "The EVE SSO Client ID")
	serverCmd.Flags().String(EveSecretKey, DefaultEmpty, "The EVE SSO Secret Key")
	serverCmd.Flags().String(EveRedirectURI, DefaultEmpty, "The EVE SSO Redirect URI")

	// TODO: this should actually be a bool but they behave wierdly
	serverCmd.Flags().String(CacheManufacturingFlag, DefaultCacheManufacturing, "Specifies, whether to regularly cache manufacturing during the runtime of the server")

	viper.BindPFlag(ListenFlag, serverCmd.Flags().Lookup(ListenFlag))
	viper.BindPFlag(RedisFlag, serverCmd.Flags().Lookup(RedisFlag))
	viper.BindPFlag(MongoFlag, serverCmd.Flags().Lookup(MongoFlag))
	viper.BindPFlag(CorporationIDFlag, serverCmd.Flags().Lookup(CorporationIDFlag))
	viper.BindPFlag(SlackAPITokenFlag, serverCmd.Flags().Lookup(SlackAPITokenFlag))
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
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.Info("Starting server...")

	if !cache.InitSSO(
		viper.GetString(EveClientID),
		viper.GetString(EveSecretKey),
		viper.GetString(EveRedirectURI)) {
		log.Errorf("Could not initialize EVE SSO, please specify a client ID, secret and redirect URI.")
		return
	}

	cache.InitCache(viper.GetString(RedisFlag))
	db.InitMongoDB(viper.GetString(MongoFlag))

	ImportSDE()

	go slack.Bot(viper.GetString(SlackAPITokenFlag))

	if viper.GetBool(CacheManufacturingFlag) {
		go ServerLoop()
	}

	router := handlers.LoggingHandler(&DebugLogWriter{Component: "http"}, routes.NewRouter(int32(viper.GetInt(CorporationIDFlag))))
	err := http.ListenAndServe(viper.GetString(ListenFlag), router)

	log.Errorf("An error occured: %v", err)
}

// ImportSDE reads the current SDE version from sde.version and imports it into the DB, if necessary.
func ImportSDE() {
	data, err := ioutil.ReadFile("sde.version")
	if err != nil {
		log.Error("Could not read SDE version, skipping import.")
		return
	}

	array := strings.Split(string(data), "-")
	if len(array) != 2 {
		log.Error("Could not read SDE version, skipping import.")
		return
	}

	i, err := strconv.Atoi(array[0])
	version := int32(i)
	server := array[1]
	if err != nil {
		log.Error("Could not read SDE version, skipping import.")
		return
	}

	log.Infof("Checking, if SDE %d is already cached...", version)

	sde := db.StaticDataExport{}
	cache.ReadCachedObject(fmt.Sprintf("sde:%d", version), &sde)

	if sde.Version == 0 {
		db.ImportSDE(version, array[1])
		sde.Version = version
		sde.Server = server

		cache.WriteCachedObject(sde)
	} else {
		log.Infof("SDE %d is already imported.", version)
	}
}

// ServerLoop takes care of reguarly caching prices and manufacturing.
func ServerLoop() {
	builderID := int32(90821267)
	builder := model.Character{}
	cache.GetCharacter(builderID, &builder)

	typeIDs := []int32{}

	var productTypeIDs []int32
	var err error
	if productTypeIDs, err = cache.GetProductTypeIDs(); err != nil {
		return
	}

	typeIDs = append(typeIDs, productTypeIDs...)
	typeIDs = append(typeIDs, db.GetMaterialTypeIDs("manufacturing")...)
	typeIDs = append(typeIDs, db.GetMaterialTypeIDs("invention")...)

	cache.GetPrices(model.JitaRegionID, typeIDs)

	for {
		// this will cache all manufacturing objects, every hour
		for _, typeID := range productTypeIDs {
			m := manufacturing.Manufacturing{}

			if err := manufacturing.NewManufacturing(builder, int32(typeID), &m); err == nil {
				cache.WriteCachedObject(m)
			}
		}

		time.Sleep(time.Duration(1) * time.Hour)
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	if err := serverCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
