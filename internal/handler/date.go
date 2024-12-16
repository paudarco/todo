package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo/internal/models"
)

// NextDate возвращает в теле ответа следующую дату по заданным параметрам.
func (h *Handler) NextDate(c *gin.Context) {
	// Получаем значения текущей даты, даты события и правила повторения из запроса
	nowStr := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")

	if len(nowStr) == 0 || len(date) == 0 || len(repeat) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "missing required parameters")
		return
	}

	now, err := time.Parse(models.Format, nowStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	nextDate, err := h.service.NextDate(now, date, repeat)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"next_date": nextDate,
	// })

	c.String(http.StatusOK, nextDate)

}
