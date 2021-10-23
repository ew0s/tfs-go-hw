package handler

import (
	"chat/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createGlobalMessage(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input entities.Message
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.Sender = id

	if err := h.services.Message.Create(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "message successfully sent to global chat",
	})
}

func (h *Handler) getGlobalMessages(c *gin.Context) {
	messages := h.services.Message.GetGlobalMessages()
	c.JSON(http.StatusOK, messages)
}
