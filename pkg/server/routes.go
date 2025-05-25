package server

import (
	"github.com/loukaspe/jedi-team-challenge/internal/core/services"
	"github.com/loukaspe/jedi-team-challenge/internal/handlers"

	"github.com/loukaspe/jedi-team-challenge/pkg/auth"
	"net/http"
	"os"
)

func (s *Server) initializeRoutes() {
	// health check
	healthCheckHandler := handlers.NewHealthCheckHandler()
	s.router.HandleFunc("/health-check", healthCheckHandler.HealthCheckController).Methods("GET")

	// auth
	jwtMechanism := auth.NewAuthMechanism(
		os.Getenv("JWT_SECRET_KEY"),
		os.Getenv("JWT_SIGNING_METHOD"),
	)
	jwtService := services.NewJwtService(jwtMechanism)
	jwtMiddleware := handlers.NewAuthenticationMw(jwtMechanism)
	jwtHandler := handlers.NewJwtClaimsHandler(jwtService, s.logger)

	s.router.HandleFunc("/token", jwtHandler.JwtTokenController).Methods(http.MethodPost)

	protected := s.router.PathPrefix("/").Subrouter()
	protected.Use(jwtMiddleware.AuthenticationMW)

	//protected.HandleFunc("/users/{user_id:[0-9]+}/favourites", createUserFavouriteHandler.AddUserFavouriteAssetController).Methods("POST")

}
