module github.com/jrapoport/gothic

go 1.17

replace github.com/containerd/containerd v1.5.5 => github.com/containerd/containerd v1.5.7

require (
	github.com/badoux/checkmail v1.2.1
	github.com/cenkalti/backoff/v4 v4.1.1
	github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804
	github.com/flashmob/go-guerrilla v1.6.1
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/cors v1.2.0
	github.com/go-chi/httprate v0.5.2
	github.com/go-chi/httptracer v0.3.0
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-playground/form/v4 v4.2.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gookit/event v1.0.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/imdario/mergo v0.3.12
	github.com/improbable-eng/go-httpwares v0.0.0-20200609095714-edc8019f93cc
	github.com/jackc/pgx/v4 v4.13.0
	github.com/joho/godotenv v1.4.0
	github.com/jrapoport/sillyname-go v0.0.0-20191016072109-82c270b69bff
	github.com/lestrrat-go/jwx v1.2.7
	github.com/lestrrat-go/test-mysqld v0.0.0-20190527004737-6c91be710371
	github.com/lib/pq v1.10.3
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/manifoldco/promptui v0.8.0
	github.com/markbates/goth v1.68.0
	github.com/matcornic/hermes/v2 v2.1.0
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/opentracing/opentracing-go v1.2.0
	github.com/orlangure/gnomock v0.18.0
	github.com/puerkitobio/purell v1.1.1
	github.com/segmentio/encoding v0.2.22
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/vcraescu/go-paginator/v2 v2.0.0
	github.com/xhit/go-simple-mail/v2 v2.10.0
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/DataDog/dd-trace-go.v1 v1.33.0
	gopkg.in/data-dog/go-sqlmock.v2 v2.0.0-20180914054222-c19298f520d0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/datatypes v1.0.2
	gorm.io/driver/clickhouse v0.2.1
	gorm.io/driver/mysql v1.1.2
	gorm.io/driver/postgres v1.1.2
	gorm.io/driver/sqlite v1.1.6
	gorm.io/driver/sqlserver v1.1.0
	gorm.io/gorm v1.21.16
)

require (
	cloud.google.com/go v0.93.3 // indirect
	github.com/ClickHouse/clickhouse-go v1.4.5 // indirect
	github.com/DataDog/datadog-go v4.4.0+incompatible // indirect
	github.com/DataDog/sketches-go v1.0.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.16.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/PuerkitoBio/goquery v1.5.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/aokoli/goutils v1.0.1 // indirect
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/containerd/containerd v1.5.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.0-20210816181553-5444fa50b93d // indirect
	github.com/denisenkom/go-mssqldb v0.11.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v20.10.8+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-chi/chi v1.5.4 // indirect
	github.com/goccy/go-json v0.7.8 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jaytaylor/html2text v0.0.0-20180606194806-57d518f124b0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a // indirect
	github.com/klauspost/cpuid/v2 v2.0.6 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/lestrrat-go/tcputil v0.0.0-20180223003554-d3c7f98154fb // indirect
	github.com/lunixbochs/vtclean v0.0.0-20180621232353-2d01aacdc34a // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/markbates/going v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/mrjones/oauth v0.0.0-20180629183705-f4e24b6d100c // indirect
	github.com/olekukonko/tablewriter v0.0.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/segmentio/asm v1.1.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	github.com/vanng822/css v0.0.0-20190504095207-a21e860bcd04 // indirect
	github.com/vanng822/go-premailer v0.0.0-20191214114701-be27abe028fe // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210924151903-3ad01bbaa167 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210925032602-92d5a993a665 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210924002016-3dee208752a0 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
