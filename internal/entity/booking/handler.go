package booking

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) Book(c *gin.Context) {
	placeID, err := strconv.Atoi(c.Query("place_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place_id"})
		return
	}
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	from, err := time.Parse(time.RFC3339, c.Query("from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
		return
	}
	to, err := time.Parse(time.RFC3339, c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
		return
	}
	if !from.Before(to) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from must be before to"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.svc.Create(ctx, userID, placeID, from, to)
	if err == ErrConflict {
		c.Status(http.StatusConflict)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) BookList(c *gin.Context) {
	userRaw := c.Query("user_id")
	placeRaw := c.Query("place_id")
	if (userRaw == "" && placeRaw == "") || (userRaw != "" && placeRaw != "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exactly one of user_id or place_id must be provided"})
		return
	}

	var userID *int
	var placeID *int
	if userRaw != "" {
		parsed, err := strconv.Atoi(userRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = &parsed
	} else {
		parsed, err := strconv.Atoi(placeRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place_id"})
			return
		}
		placeID = &parsed
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.List(ctx, userID, placeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, b := range items {
		out = append(out, gin.H{
			"id":       b.ID,
			"user_id":  b.UserID,
			"place_id": b.PlaceID,
			"from":     b.From.UTC().Format(time.RFC3339),
			"to":       b.To.UTC().Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{"bookings": out})
}
