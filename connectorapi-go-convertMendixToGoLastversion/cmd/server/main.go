package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	tcp_client_adapter "connectorapi-go/internal/adapter/client"
	handler_adapter "connectorapi-go/internal/adapter/handler/api"
	repo_adapter "connectorapi-go/internal/adapter/utils"
	service_core "connectorapi-go/internal/core/service"
	"connectorapi-go/pkg/config"
	"connectorapi-go/pkg/logger"
	"connectorapi-go/pkg/metrics"
)

// @title           Connector API Gateway
// @version         1.0
// @description     This is the API Gateway for ConnectorAPI.
// @termsOfService  http://swagger.io/terms/

// @contact.name   SYE
// @contact.url    https://aeon.co.th
// @contact.email  sye@aeon.co.th

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      10.254.97.103:8082
// @BasePath  /

// @schemes http https
func main() {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	apiKeys, err := config.LoadAPIKeys("./configs/apikeys.json")
	if err != nil {
		log.Fatalf("FATAL: Failed to load apiKeys: %v", err)
	}

	dr, err := config.LoadDestinationsAndRoutes("./configs/destinations_routes.json")
	if err != nil {
		log.Fatalf("FATAL: Failed to load destinations and routes: %v", err)
	}

	appLogger := logger.New(cfg.Logger.Level)
	defer appLogger.Sync()
	appLogger.Info("Logger initialized")

	gin.SetMode(cfg.Server.Mode)
	appLogger.Infow("Gin mode set", "mode", cfg.Server.Mode)

	appLogger.Info("Initializing dependencies...")
	metrics.Init()

	// --- Adapters ---
	apiKeyRepo := repo_adapter.NewAPIKeyRepository(apiKeys)

	// --- TCP Socket Client Initialization ---
	tcpClient := tcp_client_adapter.NewBasicTCPSocketClient(
		5*time.Second,  // Dial Timeout (e.g., 5 seconds to establish connection)
		10*time.Second, // Read/Write Timeout (e.g., 10 seconds for data transfer)
	)
	appLogger.Info("TCP Socket Client initialized")
	
	appLogger.Infow("Loaded routes", "routes", dr.Routes)
	appLogger.Infow("Loaded destinations", "destinations", dr.Destinations)

	// --- Core Services ---
	collectionService := service_core.NewCollectionService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	agreementService := service_core.NewAgreementService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	creditcardService := service_core.NewCreditCardService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	commonService := service_core.NewCommonService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	selfServiceService := service_core.NewSelfServiceService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	registerService := service_core.NewRegisterService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	customerLowerService := service_core.NewCustomerLowerService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	consentService := service_core.NewConsentService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	uhpService := service_core.NewUhpService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	mobileService := service_core.NewMobileService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	applicationCapService := service_core.NewApplicationCapService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	applicationLowerService := service_core.NewApplicationLowerService(cfg, appLogger, tcpClient, dr.Routes, dr.Destinations)
	appLogger.Info("Customer Service initialized with TCP client")

	// --- Handlers (API Layer) ---
	collectionHandler := handler_adapter.NewCollectionHandler(collectionService, appLogger, apiKeyRepo, cfg)
	agreementHandler := handler_adapter.NewAgreementHandler(agreementService, appLogger, apiKeyRepo, cfg)
	creditcardHandler := handler_adapter.NewCreditCardHandler(creditcardService, appLogger, apiKeyRepo, cfg)
	commonHandler := handler_adapter.NewCommonHandler(commonService, appLogger, apiKeyRepo, cfg)
	selfServiceHandler := handler_adapter.NewSelfServiceHandler(selfServiceService, appLogger, apiKeyRepo, cfg)
	registerHandler := handler_adapter.NewRegisterHandler(registerService, appLogger, apiKeyRepo, cfg)
	customerLowerHandler := handler_adapter.NewCustomerLowerHandler(customerLowerService, appLogger, apiKeyRepo, cfg)
	consentHandler := handler_adapter.NewConsentHandler(consentService, appLogger, apiKeyRepo, cfg)
	uhpHandler := handler_adapter.NewUhpHandler(uhpService, appLogger, apiKeyRepo, cfg)
	mobileHandler := handler_adapter.NewMobileHandler(mobileService, appLogger, apiKeyRepo, cfg)
	applicationCapHandler := handler_adapter.NewApplicationCapHandler(applicationCapService, appLogger, apiKeyRepo, cfg)
	applicationLowerHandler := handler_adapter.NewApplicationLowerHandler(applicationLowerService, appLogger, apiKeyRepo, cfg)

	appLogger.Info("Setting up router...")
	router := handler_adapter.SetupRouter(appLogger, apiKeyRepo, collectionHandler, agreementHandler, creditcardHandler, commonHandler, selfServiceHandler, registerHandler, customerLowerHandler, consentHandler, uhpHandler, mobileHandler, applicationCapHandler, applicationLowerHandler)

	serverAddress := fmt.Sprintf(":%s", cfg.Server.Port)
	appLogger.Infow("Starting server", "address", serverAddress)
	if err := router.Run(serverAddress); err != nil {
		appLogger.Fatalw("Failed to start server", "error", err)
	}
}
