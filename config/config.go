package config

import (
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var Set = wire.NewSet(NewConfig)

type Config struct {
	Server         Server
	DB             DB
	Cache          Cache
	AdapterService AdapterService
	RabbitMQ       RabbitMQ
	CronJob        CronJob
	Notify         Notify
}

type Server struct {
	KeyID               string
	Name                string
	ApiHeaderKey        string
	AppVersion          string
	RestPort            string
	GrpcPort            string
	BaseURI             string
	Mode                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	SSL                 bool
	CtxDefaultTimeout   int
	CSRF                bool
	Debug               bool
	MaxCountRequest     int           // max count of connections
	ExpirationLimitTime time.Duration //  expiration time of the limit
}

type Cache struct {
	Redis Redis
}

type Redis struct {
	Address  string
	Port     int
	Password string
	DB       int
}
type CronJob struct {
	CheckingTxStatus string
}

type DB struct {
	Mongodb Mongodb
}

type Notify struct {
	Token        string
	GroupAlertID int64
}

type Mongodb struct {
	DbName          string
	Username        string
	Password        string
	Connection      string
	ConnectTimeout  time.Duration
	MaxConnIdleTime int
	MinPoolSize     uint64
	MaxPoolSize     uint64
}

type RabbitMQ struct {
	SagaOrderEvent          SagaOrderEvent
	SagaOrderProductEvent   SagaOrderProductEvent
	SagaOrderPromotionEvent SagaOrderPromotionEvent
	SagaOrderPaymentEvent   SagaOrderPaymentEvent
	SagaOrderDeliveryEvent  SagaOrderDeliveryEvent
	SagaOrderEmailEvent     SagaOrderEmailEvent
	SagaOrderCartEvent      SagaOrderCartEvent

	Connection  string
	ServiceName string

	OrderTransactionExchange string
	OrderCommitRoutingKey    string
}

type SagaOrderEvent struct {
	Connection        string
	Exchange          string
	PublishRoutingKey string
	CancelRoutingKey  string
	ReplyRoutingKey   string
	Queue             string
}

type SagaOrderProductEvent struct {
	Connection         string
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
	Queue              string
}

type SagaOrderPromotionEvent struct {
	Connection         string
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
	Queue              string
}

type SagaOrderPaymentEvent struct {
	Connection         string
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
	Queue              string
}

type SagaOrderDeliveryEvent struct {
	Connection         string
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
	Queue              string
}

type SagaOrderEmailEvent struct {
	Connection         string
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
	Queue              string
}

type SagaOrderCartEvent struct {
	Connection         string
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
	Queue              string
}

type AdapterService struct {
	AuthService      AuthService
	UserService      UserService
	ProductService   ProductService
	EmailService     EmailService
	DeliveryService  DeliveryService
	PromotionService PromotionService
	StoreService     StoreService
}

type AuthService struct {
	BaseURL     string
	InternalKey string
}

type UserService struct {
	UserURL     string
	InternalKey string
}

type ProductService struct {
	BaseURL     string
	InternalKey string
}

type StoreService struct {
	BaseURL     string
	InternalKey string
}

type EmailService struct {
	Email string
	Host  string
	Key   string
}

type DeliveryService struct {
	BaseURL     string
	InternalKey string
}

type PromotionService struct {
	BaseURL     string
	InternalKey string
}

// Get config path for local or docker
func getDefaultConfig() string {
	return "/config/config"
}

// Load config file from given path
func NewConfig() (*Config, error) {
	config := Config{}
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}
	fmt.Printf("config path:%s\n", path)

	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}
