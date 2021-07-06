module github.com/filecoin-project/venus-auth

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/filecoin-project/go-address v0.0.5
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gbrlsnchs/jwt/v3 v3.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-resty/resty/v2 v2.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/google/uuid v1.2.0
	github.com/influxdata/influxdb-client-go/v2 v2.2.2
	github.com/ipfs-force-community/metrics v0.0.0-20210705093944-918711d7932a
	github.com/jmoiron/sqlx v1.3.3
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/magiconair/properties v1.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cast v1.3.0
	github.com/spf13/viper v1.3.2
	github.com/ugorji/go v1.2.4 // indirect
	github.com/urfave/cli/v2 v2.3.0
	go.opencensus.io v0.23.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools v2.2.0+incompatible
)

replace github.com/google/flatbuffers => github.com/google/flatbuffers v1.12.1
