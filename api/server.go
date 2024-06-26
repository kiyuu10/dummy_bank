package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kiyuu10/simplebank/db/sqlc"
	"github.com/kiyuu10/simplebank/token"
	"github.com/kiyuu10/simplebank/util"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupRouter()

	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	router.POST("/users/create", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/token/refresh", server.RenewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.POST("/accounts/get", server.getAccount)
	authRoutes.POST("/accounts/list", server.listAccount)
	authRoutes.POST("/accounts/list-by-owner", server.listAccountByOwner)
	authRoutes.POST("/transfers/create", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
