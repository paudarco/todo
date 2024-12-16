package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo/internal/errors"
	"github.com/paudarco/todo/internal/models"
)

// SignIn возвращает JWT-токен для входа в систему.
func (h *Handler) SignIn(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, errors.ErrInvalidPassword.Error())
		return
	}

	response, err := h.service.SignIn(user)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, errors.ErrInvalidPassword.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": response.Token,
	})

}
