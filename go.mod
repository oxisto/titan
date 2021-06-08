module github.com/oxisto/titan

go 1.14

require (
	github.com/antihax/goesi v0.0.0-20210507032439-ea5f011d207e
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/structs v1.1.0
	github.com/gin-contrib/static v0.0.1
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/oxisto/bellows v1.0.0
	github.com/oxisto/evesso v1.0.6
	github.com/oxisto/go-httputil v1.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
