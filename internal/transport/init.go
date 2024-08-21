package transport

import (
	"net/http"
	"time"

	_ "github.com/amiosamu/gofemart/docs"
	"github.com/amiosamu/gofemart/internal/config"
	"github.com/amiosamu/gofemart/internal/hash"
	"github.com/amiosamu/gofemart/internal/repository"
	"github.com/amiosamu/gofemart/internal/service"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type APIServer struct {
	config        *config.Config
	router        *chi.Mux
	logger        *log.Logger
	storage       *repository.Storage
	users         *service.Users
	orders        *service.Orders
	withdraw      *service.Bonuses
	scoringsystem *service.ScoringSystem
}

func NewAPIServer(config *config.Config) *APIServer {
	return &APIServer{
		config: config,
		router: chi.NewRouter(),
		logger: log.New(),
	}
}

func (s *APIServer) Start() error {
	s.config.ParseFlags()
	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return err
	}

	db, err := s.configureStore()
	if err != nil {
		return err
	}
	s.storage = db
	defer s.storage.Close()

	hasher := hash.NewSHA1Hasher("salt")
	s.users = service.NewUsers(db, hasher, []byte("sample secret"), s.config.TokenTTL)
	s.orders = service.NewOrders(db)
	s.withdraw = service.NewBonuses(db, db)
	s.scoringsystem = service.NewScoringSystem(db)

	s.logger.Info("starting api server")

	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.ScoringSystem()
			default:
				continue
			}

		}
	}()

	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(withLogging)
	s.router.Post("/api/user/register", s.SighUp)
	s.router.Post("/api/user/login", s.SighIn)
	s.router.With(s.authMiddleware).Post("/api/user/orders", s.OrderUploading)
	s.router.With(s.authMiddleware).Get("/api/user/orders", s.GetAllOrders)
	s.router.With(s.authMiddleware).Get("/api/user/balance", s.Balance)
	s.router.With(s.authMiddleware).Post("/api/user/balance/withdraw", s.Withdraw)
	s.router.With(s.authMiddleware).Get("/api/user/withdrawals", s.Withdrawals)
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}

func (s *APIServer) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureStore() (*repository.Storage, error) {
	db, err := repository.NewStorage(s.config.DBPort)
	if err != nil {
		return nil, err
	}
	return db, nil
}
