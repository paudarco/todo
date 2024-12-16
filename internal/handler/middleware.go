package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo/internal/errors"
)

// AuthMiddleware проверяет авторизацию пользователя.
func (h *Handler) AuthMiddleware(c *gin.Context) {
	if len(h.cfg.Password) == 0 {
		c.Next()
		return
	}

	token, err := c.Cookie("token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	err = h.service.ValidateToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Next()
}
