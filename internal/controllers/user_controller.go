package controllers

import (
	"net/http"
	"strconv"
	"autograder/internal/services"
	"autograder/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	users *services.UserService
}

func NewUserController(users *services.UserService) *UserController {
	return &UserController{users: users}
}

func (h *UserController) ListStudents(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.users.ListStudents(c.Request.Context(), search, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	
	utils.Success(c, http.StatusOK, "students retrieved", gin.H{
		"users": users,
		"total": total,
		"page": page,
		"limit": limit,
	})
}
