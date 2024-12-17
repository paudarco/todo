package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo/internal/errors"
)

// AuthMiddleware проверяет авторизацию пользователя.
func (h *Handler) AuthMiddleware(c *gin.Context) {
	// Если пароль не установлен - пропускаем.
	if len(h.cfg.Password) == 0 {
		c.Next()
		return
	}

	// Проверяем наличие токена и его валидность.
	token, err := c.Cookie("token")
	if err != nil || len(token) == 0 {
		c.Redirect(http.StatusSeeOther, "/login")
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	err = h.service.ValidateToken(token)
	if err != nil {
		// Если токен невалиден, удаляем его из куков и перенаправляем на страницу логина
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/login")
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Next()
}

// RedirectIfAuthenticated проверяет авторизацию пользователя и
// перенаправляет его на "/" (index.html)" если он авторизован.
func (h *Handler) RedirectIfAuthenticated(c *gin.Context) {
	// Проверяем наличие токена в куках
	token, err := c.Cookie("token")
	if err == nil && len(token) > 0 {
		// Проверяем валидность токена
		if err := h.service.ValidateToken(token); err == nil {
			// Если токен валиден, перенаправляем на "/"
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}
	}

	// Если токена нет или он невалиден, продолжаем обработку
	c.Next()
}
