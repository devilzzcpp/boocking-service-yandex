package app

import (
	"booking_service/internal/entity/booking"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(lg *zap.Logger, h *booking.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestLogger(lg))

	r.GET("/ping", h.Ping)
	r.POST("/book", h.Book)
	r.GET("/booklist", h.BookList)

	return r
}
