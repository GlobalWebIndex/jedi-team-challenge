package server

import (
	"github.com/loukaspe/jedi-team-challenge/internal/core/services"
	"github.com/loukaspe/jedi-team-challenge/internal/handlers"
	"github.com/loukaspe/jedi-team-challenge/internal/handlers/chatSessions"
	"github.com/loukaspe/jedi-team-challenge/internal/repositories"

	"github.com/loukaspe/jedi-team-challenge/pkg/auth"
	"net/http"
	"os"
)

//	@title			Louk Chatwalker
//	@version		1.0
//	@description	GWI's Jedi Team Challenge

//	@host		localhost:8080
//	@BasePath	/

//	@contact.name	Loukas Peteinaris
//	@contact.url	loukas.peteinaris@gmail.com

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Header value should be in the form of `Bearer <JWT access token>`

// @accept		json
// @produce	json
func (s *Server) initializeRoutes() {
	// health check
	healthCheckHandler := handlers.NewHealthCheckHandler(s.DB)
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

	chatSessionRepository := repositories.NewChatSessionRepository(s.DB)
	chatSessionService := services.NewChatSessionService(s.logger, chatSessionRepository)

	createChatSessionHandler := chatSessions.NewCreateUserChatSessionHandler(chatSessionService, s.logger)

	protected.HandleFunc("/users/{user_id}/chat-sessions", createChatSessionHandler.CreateUserChatSessionAssetController).Methods("POST")

}
