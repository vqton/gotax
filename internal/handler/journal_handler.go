package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
)

func (h *Handler) CreateEntry(c *gin.Context) {
	var entry domain.JournalEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.CreateEntry(c.Request.Context(), &entry, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *Handler) SubmitEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.SubmitForReview(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted for review"})
}

func (h *Handler) ReviewEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.ReviewEntry(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry reviewed"})
}

func (h *Handler) ApproveEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.ApproveEntry(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry approved"})
}

func (h *Handler) PostJournalEntry(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostEntry(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry posted"})
}

func (h *Handler) CancelJournalEntry(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CancelEntry(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry cancelled"})
}

func (h *Handler) GetJournalEntry(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.svc.GetEntryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *Handler) GetJournalEntries(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	status := c.Query("status")
	if fromStr != "" && toStr != "" {
		from, err1 := time.Parse("2006-01-02", fromStr)
		to, err2 := time.Parse("2006-01-02", toStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		entries, err := h.svc.GetEntriesByDateRange(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entries)
		return
	}
	if status != "" {
		entries, err := h.svc.GetEntriesByStatus(c.Request.Context(), domain.JournalEntryStatus(status))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entries)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "specify from/to dates or status"})
}
