package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getTodo(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	publicKey, privateKey, err := getUserAPIKeys(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":     userID,
		"public_key":  publicKey,
		"private_key": privateKey,
		"todo":        "some todo",
	})
}
