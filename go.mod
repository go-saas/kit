module github.com/goxiaoy/go-saas-kit

go 1.18

require (
	github.com/alexedwards/argon2id v0.0.0-20211130144151-3585854a6387
	github.com/casbin/casbin/v2 v2.42.0
	github.com/casbin/gorm-adapter/v3 v3.5.0
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-kratos/kratos/v2 v2.2.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/handlers v1.5.1
	github.com/goxiaoy/go-eventbus v0.0.5
	github.com/goxiaoy/go-saas v0.1.1-0.20220402162851-7116458f7dea
	github.com/goxiaoy/gorm-concurrency v1.0.5
	github.com/goxiaoy/uow v0.0.0-20210815151702-b0032203778a
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.8.0
	github.com/mennanov/fmutils v0.1.1
	github.com/mitchellh/mapstructure v1.4.3
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354
	github.com/stretchr/testify v1.7.1
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220315194320-039c03cc5b86 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/mysql v1.3.3-0.20220223060847-68a051989967
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.4
)

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/Shopify/sarama v1.32.0
	github.com/aws/aws-sdk-go v1.43.19
	github.com/fclairamb/afero-s3 v0.3.1
	github.com/go-redis/cache/v8 v8.4.3
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.5
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/goxiaoy/sessions v1.2.2-0.20220327125603-baf0b311512e
	github.com/nicksnyder/go-i18n/v2 v2.2.0
	github.com/nyaruka/phonenumbers v1.0.74
	github.com/ory/hydra-client-go v1.11.7
	github.com/samber/lo v1.10.1
	github.com/spf13/afero v1.8.2
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.1.10
	github.com/xhit/go-simple-mail/v2 v2.11.0
	go.opentelemetry.io/otel v1.6.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.0
	go.opentelemetry.io/otel/sdk v1.6.0
	go.opentelemetry.io/otel/trace v1.6.0
	golang.org/x/text v0.3.7
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.12.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-redis/redis/extra/rediscmd/v8 v8.11.5 // indirect
	github.com/go-test/deep v1.0.8 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.0.0-20170517235910-f1bb20e5a188 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/jackc/pgx/v4 v4.15.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/toorop/go-dkim v0.0.0-20201103131630-e1cd1a0a5208 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.1.10 // indirect
	github.com/vmihailenco/go-tinylfu v0.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.4 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.6.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/metric v0.27.0 // indirect
	go.opentelemetry.io/proto/otlp v0.12.0 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/exp v0.0.0-20220314205449-43aec2f8a4e7 // indirect
	golang.org/x/oauth2 v0.0.0-20220309155454-6242fa91716a // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gorm.io/driver/postgres v1.3.1 // indirect
	gorm.io/driver/sqlserver v1.3.1 // indirect
	gorm.io/plugin/dbresolver v1.1.0 // indirect
)

replace github.com/casbin/gorm-adapter/v3 => github.com/Goxiaoy/gorm-adapter/v3 v3.5.1-0.20220326110105-403c86d95e88

replace github.com/go-redis/cache/v8 => github.com/Goxiaoy/cache/v8 v8.4.4-0.20220418161131-a691bbf092a2
