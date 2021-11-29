package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/pkg/krakenFuturesSDK"
)

// @Summary SendOrder
// @Security ApiKeyAuth
// @Tags orderManager
// @Description sendOrder to kraken futures API
// @ID sendOrder
// @Accept  json
// @Produce  json
// @Param input body krakenFuturesSDK.SendOrderArguments true "send order info"
// @Success 200 {string} string "order_id"
// @Failure 400,401,404 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure default {object} errResponse
// @Router /orderManager/send-order [post]
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

func (h *Handler) startTrade(c *gin.Context) {
	conn, err := h.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	defer conn.Close()

	userID, err := getUserID(c)
	if err != nil {
		newWebsocketErrResponse(c, http.StatusUnauthorized, conn, err.Error())
		return
	}

	var td types.TradingDetails
	if err := conn.ReadJSON(&td); err != nil {
		newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
		return
	}

	if err := h.validate.Struct(td); err != nil {
		newWebsocketErrResponse(c, http.StatusBadRequest, conn, err.Error())
		return
	}

	orderID, err := h.services.KrakenOrdersManager.StartTrading(userID, td)
	if err != nil {
		newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
		return
	}

	output := map[string]interface{}{
		"order_id": orderID,
	}
	if err := conn.WriteJSON(output); err != nil {
		newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
		return
	}
}
