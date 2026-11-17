package teacher

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// NotificationHandler handles HTTP requests for notification endpoints.
type NotificationHandler struct {
	svc NotificationService
}

// NewNotificationHandler constructs a new NotificationHandler instance.
func NewNotificationHandler(svc NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// ListNotifications handles GET /api/notifications.
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		handler.RespondError(c, handler.ErrUnauthorized)
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	notifs, meta, err := h.svc.GetNotifications(c.Request.Context(), userID, page, limit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    notifs,
		"meta":    meta,
	})
}

// GetUnreadCount handles GET /api/notifications/unread-count.
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		handler.RespondError(c, handler.ErrUnauthorized)
		return
	}

	res, err := h.svc.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// MarkAsRead handles PATCH /api/notifications/:id/read.
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		handler.RespondError(c, handler.ErrUnauthorized)
		return
	}

	idStr := strings.TrimSpace(c.Param("id"))
	notifID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID notifikasi tidak valid"))
		return
	}

	if err := h.svc.MarkAsRead(c.Request.Context(), userID, notifID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi ditandai sudah dibaca",
	})
}

// MarkAllAsRead handles PATCH /api/notifications/read-all.
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		handler.RespondError(c, handler.ErrUnauthorized)
		return
	}

	if err := h.svc.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Semua notifikasi ditandai sudah dibaca",
	})
}
