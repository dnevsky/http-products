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
		newErrorResponseJSON(c, http.StatusBadRequest, map[string]interface{}{
			"message":  "bad data for limit",
			"rawlimit": rawLimit,
			"limit":    limit,
		})
		return
	}

	offset, err := strconv.Atoi(rawOffset)
	if err != nil {
		h.logger.Infow("bad data for offset",
			"rawoffser", rawOffset,
			"offset", offset,
		)
		newErrorResponseJSON(c, http.StatusBadRequest, map[string]interface{}{
			"message":   "bad data for offset",
			"rawoffset": rawOffset,
			"offset":    offset,
		})
		return
	}

	if offset < 0 || limit <= 0 {
		h.logger.Infow("invalid offset or limit",
			"offset", offset,
			"limit", limit,
		)
		newErrorResponseJSON(c, http.StatusBadRequest, map[string]interface{}{
			"message": "invalid offset or limit",
			"offset":  offset,
			"limit":   limit,
		})
		return
	}

	products, err := h.services.Product.GetAll(c, limit, offset)
	if err != nil {
		h.logger.Infow("error while get all products",
			"products", products, "error", err.Error(),
		)
		newErrorResponse(c, http.StatusInternalServerError, "error while get all products")
		return
	}

	c.JSON(http.StatusOK, products)
}
