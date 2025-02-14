package missions

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"spy-cats/internal/models/missions"
	missionService "spy-cats/internal/service/missions"
	"strconv"
)

type TargetHandler struct {
	targetService missionService.Targets
	logger        *slog.Logger
}

func NewTargetHandler(targetService missionService.Targets, logger *slog.Logger) *TargetHandler {
	return &TargetHandler{
		targetService: targetService,
		logger:        logger,
	}
}

func (h *TargetHandler) Create(c *gin.Context) {
	var target missions.Target
	err := c.BindJSON(&target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target JSON"})
		return
	}

	if target.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target name"})
		return
	}

	if target.Country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target country"})
		return
	}

	if target.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target cannot be completed from the beginning"})
		return
	}

	if err = h.targetService.Create(target); err != nil {
		h.logger.Error("failed to create target: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *TargetHandler) AddTargetToMission(c *gin.Context) {
	targetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse target id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect target id"})
		return
	}

	missionID, err := strconv.Atoi(c.Param("mission_id"))
	if err != nil {
		h.logger.Error("incorrect mission id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect mission id"})
		return
	}

	if err := h.targetService.AddTargetToMission(targetId, missionID); err != nil {
		h.logger.Error("failed to add target to mission: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add target to mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *TargetHandler) DeleteTargetFromMission(c *gin.Context) {
	targetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse target id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect target id"})
		return
	}

	if err := h.targetService.DeleteTargetFromMission(targetId); err != nil {
		h.logger.Error("failed to delete target from mission: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete target from mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *TargetHandler) CompleteTarget(c *gin.Context) {
	targetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect target id"})
		return
	}

	if err := h.targetService.CompleteTarget(targetId); err != nil {
		h.logger.Error("failed to complete target: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *TargetHandler) UpdateNotes(c *gin.Context) {
	targetId, err := strconv.Atoi(c.Param("id"))
	if targetId < 0 || err != nil {
		h.logger.Error("failed to parse target id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	var input struct {
		Notes string `json:"notes"`
	}

	if err = c.BindJSON(&input); err != nil {
		h.logger.Error("failed to parse target notes: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if err = h.targetService.UpdateNotes(targetId, input.Notes); err != nil {
		h.logger.Error("failed to update notes: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
