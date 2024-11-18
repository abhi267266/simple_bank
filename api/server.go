package api

import (
	"github.com/gin-gonic/gin"
	db "tutorial.sqlc.dev/app/db/sqlc"
)

type Server struct {
	store *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
// It takes a db.Store as an argument, initializes a gin router,
// and registers the createAccount handler under the "/accounts" endpoint.
// The function returns a pointer to the created Server.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)

    server.router = router
    return server
}

func errorResponse(err error) gin.H {
    return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
    return server.router.Run(address)
}



