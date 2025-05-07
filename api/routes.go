package api

import (
	"net/http"
	"roommate-finder/campay"
	"roommate-finder/db/repo"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	querier      repo.Querier
	campayClient *campay.Requests
}

func NewControllerHandler(querier repo.Querier, campayClient *campay.Requests) *UserHandler {
	return &UserHandler{
		querier:      querier,
		campayClient: campayClient,
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

	return r
}
