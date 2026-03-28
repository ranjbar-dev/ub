package di

import (
	"exchange-go/config"
	"exchange-go/internal/communication"
	"exchange-go/internal/platform"

	"github.com/sarulabs/di"
)

// DI registrations for infrastructure and platform services.
// These are registered first as they have no internal service dependencies.
// All other services depend on one or more of: configService, loggerService,
// dbClient, redisClient, cacheService, mqttClient, rabbitmqClient, httpClient, wsClient.
func addConfigService() {
	mustAdd(di.Def{
		Name:  ConfigService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			viper := config.SetConfigs()
			return platform.NewConfigs(viper), nil
		},
	})
}

func addCacheService() {
	mustAdd(di.Def{
		Name:  cacheService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return platform.NewCache(configsService, logger), nil
		},
	})
}

func addDBClient() {
	mustAdd(di.Def{
		Name:  dbClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			return platform.NewDbClient(configsService)
		},
	})
}

func addLogger() {
	mustAdd(di.Def{
		Name:  LoggerService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			return platform.NewLogger(configsService), nil
		},
	})
}

func addWSClient() {
	mustAdd(di.Def{
		Name:  wsClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return platform.NewWsClient(), nil
		},
	})
}

func addMQTTClient() {
	mustAdd(di.Def{
		Name:  mqttClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			loggerService := ctn.Get(LoggerService).(platform.Logger)
			return platform.NewMqttClient(configsService, loggerService), nil
		},
	})
}

func addRedisClient() {
	mustAdd(di.Def{
		Name:  RedisClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			return platform.NewRedisClient(configsService), nil
		},
	})
}

func addHTTPClient() {
	mustAdd(di.Def{
		Name:  httpClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return platform.NewHTTPClient(), nil
		},
	})
}

func addMQTTManager() {
	mustAdd(di.Def{
		Name:  mqttManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			mqttCli := ctn.Get(mqttClient).(platform.MqttClient)
			return communication.NewMqttManager(mqttCli), nil
		},
	})
}

func addRabbitmqClient() {
	mustAdd(di.Def{
		Name:  rabbitmqClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return platform.NewRabbitMqClient(configService, logger), nil
		},
	})
}

func addQueueManager() {
	mustAdd(di.Def{
		Name:  queueManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			rabbitmqCli := ctn.Get(rabbitmqClient).(platform.RabbitMqClient)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return communication.NewQueueManager(rabbitmqCli, logger), nil
		},
	})
}

func addJWTHandler() {
	mustAdd(di.Def{
		Name:  jwtHandler,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return platform.NewJwtHandler(), nil
		},
	})
}

func addPasswordEncoder() {
	mustAdd(di.Def{
		Name:  passwordEncoder,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return platform.NewPasswordEncoder(), nil
		},
	})
}
