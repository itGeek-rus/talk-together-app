package main

import (
	talk_together_app "github.com/firstproject/talk-together-app/hub"
	"github.com/firstproject/talk-together-app/pkg/handler"
	"github.com/firstproject/talk-together-app/pkg/kafka"
	"github.com/firstproject/talk-together-app/pkg/redis"
	"github.com/firstproject/talk-together-app/pkg/repository"
	"github.com/firstproject/talk-together-app/pkg/service"
	"github.com/firstproject/talk-together-app/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// @title Talk together app API
// @version 1.0
// @description API Server for TalkTogether Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Error initializing DB: %s", err.Error())
	}

	redisClient := redis.NewRedisClient(
		viper.GetString("redis.addr"),
		os.Getenv("REDIS_PASSWORD"),
		viper.GetInt("redis.db"),
	)

	kafkaProducer, err := kafka.NewKafkaProducer(
		[]string{viper.GetString("kafka.brokers")},
		viper.GetString("kafka.topic"),
	)
	if err != nil {
		logrus.Fatalf("Error initializing Kafka producer: %s", err.Error())
	}

	hub := talk_together_app.NewHub()
	go hub.Run()

	kafkaConsumer, err := kafka.NewKafkaConsumer(
		[]string{viper.GetString("kafka.brokers")},
		viper.GetString("kafka.topic"),
		hub,
	)
	go kafkaConsumer.Start()

	repos := repository.NewRepository(db)
	services := service.NewService(repos, redisClient, kafkaProducer)
	handlers := handler.NewHandler(services, hub)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error starting server: %s", err.Error())
	}

	//services repos, redisClient, kafkaProducer
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
