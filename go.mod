module github.com/goxiaoy/go-saas-kit

go 1.17

require (
	entgo.io/ent v0.9.0
	github.com/a8m/rql v1.3.1-0.20210621074553-3a40179141a1
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/alexedwards/argon2id v0.0.0-20210511081203-7d35d68092b8
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/go-kratos/kratos/v2 v2.1.1
	github.com/go-kratos/swagger-api v1.0.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/goxiaoy/go-saas v0.0.10-0.20211112043352-57ebc170b313
	github.com/goxiaoy/gorm-concurrency v1.0.2
	github.com/goxiaoy/uow v0.0.0-20210815151702-b0032203778a
	github.com/mennanov/fmutils v0.1.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354
	github.com/ory/keto/proto/ory/keto/acl/v1alpha1 v0.0.0-20210616104402-80e043246cf9
	github.com/segmentio/kafka-go v0.4.18
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20211117180635-dee7805ff2e1 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.1
	gorm.io/driver/mysql v1.1.2
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.13
)

replace github.com/a8m/rql v1.3.1-0.20210621074553-3a40179141a1 => github.com/Goxiaoy/rql v1.3.1-0.20210823140701-2d9807375ca8
