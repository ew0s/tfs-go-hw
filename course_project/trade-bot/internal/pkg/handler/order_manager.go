package handler

import (
	"net/http"
	"trade-bot/pkg/krakenFuturesSDK"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sendOrder(c *gin.Context) {
	var input krakenFuturesSDK.SendOrderArguments

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	orderID, err := h.services.KrakenOrdersManager.SendOrder(userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"order_id": orderID,
	})
}

func (h *Handler) editOrder(c *gin.Context) {

}

func (h *Handler) cancelOrder(c *gin.Context) {

}

func (h *Handler) cancelAllOrders(c *gin.Context) {

}
