package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kaviraj-j/go-bank/db/sqlc"
)

// Server serves http Requests for bank application
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Creates, setup routes and returns a server
func NewServer(store db.Store) *Server {
	server := &Server{store: store}

	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)
	server.router = router

	return server
}

// Starts and runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Returns formatted error responses
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
