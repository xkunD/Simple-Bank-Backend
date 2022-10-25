package api

import (
	"fmt"
	db "go-simple-bank/db/sqlc"
	"go-simple-bank/token"
	"go-simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	// Initialize token maker and assign store
	token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create Paseto token maker: %v", err)
	}
	server := &Server{store: store, tokenMaker: token}

	// Add currency validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// Start runs the HTTP server on specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errorResponse returns an error as gin.H (JSON)
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// setupRouter wires up the routes for the server
func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.POST("/tokens/renew_access", s.renewAccessToken)

	// Authenticated routes
	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccounts)
	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router

}
