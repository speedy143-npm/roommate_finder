package api

import (
	"log"
	"net/http"
	"roommate-finder/campay"
	"roommate-finder/db/repo"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	querier      repo.Querier
	campayClient *campay.Requests
}

func NewControllerHandler(querier repo.Querier, campayClient *campay.Requests) *UserHandler {
	handler := &UserHandler{
		querier:      querier,
		campayClient: campayClient,
	}
	go handler.startCleanupJob()
	return handler
}

// goroutine to automatically cleanup the database by deleting expired resetTokens
func (h *UserHandler) startCleanupJob() {
	ticker := time.NewTicker(4 * time.Hour) // Runs every 10 minutes
	defer ticker.Stop()

	for range ticker.C {
		err := h.querier.DeleteExpiredTokens(&gin.Context{}) // Adjust context as needed
		if err != nil {
			log.Println("Error deleting expired tokens:", err)
		} else {
			log.Println("Expired tokens cleaned up successfully")
		}
	}
}

func (h *UserHandler) WireHttpHandler() http.Handler {
	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// add routes
	r.POST("/register", h.handleUserRegistration)
	r.POST("/match/:id1/:id2", h.handleUserMatch)
	r.GET("/user/:id", h.handleGetUser)
	r.PATCH("/user/:id", h.handleUpdateUser)
	r.PATCH("/forgot-password", h.handleForgotPassword)
	r.PATCH("/reset-password", h.handleResetPassword)

	return r
}
