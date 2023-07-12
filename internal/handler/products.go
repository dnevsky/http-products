package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAll(c *gin.Context) {
	rawLimit := c.Query("limit")
	rawOffset := c.Query("offset")

	// проверяем, что-бы это были числа
	// проверки на пустые limit & offset можно не делать, strconv всё сделал, strconv умный
	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		h.logger.Infow("bad data for limit",
			"rawlimit", rawLimit,
			"limit", limit,
		)
		newErrorResponse(c, http.StatusBadRequest, "bad data for limit")
		return
	}

	offset, err := strconv.Atoi(rawOffset)
	if err != nil {
		h.logger.Infow("bad data for offset",
			"rawoffser", rawOffset,
			"offset", offset,
		)
		newErrorResponse(c, http.StatusBadRequest, "bad data for offset")
		return
	}

	products, err := h.services.Product.GetAll(limit, offset)
	if err != nil {
		h.logger.Infow("error while get all products",
			"products", products, "error", err.Error(),
		)
		newErrorResponse(c, http.StatusInternalServerError, "error while get all products")
		return
	}

	c.JSON(http.StatusOK, products)
}
