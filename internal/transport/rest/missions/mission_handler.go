package missions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"spy-cats/internal/models/missions"
	missionService "spy-cats/internal/service/missions"
	"strconv"
)

type Handler struct {
	missionService missionService.Missions
	logger         *slog.Logger
}

func NewHandler(missionService missionService.Missions, logger *slog.Logger) *Handler {
	return &Handler{
		missionService: missionService,
		logger:         logger,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	allMissions, err := h.missionService.GetAll()
	if err != nil {
		h.logger.Error("failed to get all missions: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all missions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": allMissions})
}

func (h *Handler) Create(c *gin.Context) {
	var mission missions.Mission
	err := c.BindJSON(&mission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission JSON"})
		return
	}

	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid targets quantity"})
		return
	}
	if mission.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mission cannot be completed from the beginning"})
		return
	}

	if err = h.missionService.Create(mission); err != nil {
		h.logger.Error("failed to create mission: ", err)
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *Handler) GetById(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}
	mission, err := h.missionService.GetById(missionID)
	if errors.Is(err, missions.ErrMissionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		h.logger.Error("failed to get mission: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": mission})
}

func (h *Handler) Delete(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}

	err = h.missionService.Delete(missionID)
	if err != nil {
		h.logger.Error("failed to delete mission: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *Handler) AssignCat(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect mission id"})
		return
	}
	catId, err := strconv.Atoi(c.Param("cat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect cat id"})
		return
	}

	if err := h.missionService.AssignCat(missionID, catId); err != nil {
		h.logger.Error("failed to assign cat: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *Handler) CompleteMission(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect mission id"})
		return
	}

	if err := h.missionService.CompleteMission(missionID); err != nil {
		h.logger.Error("failed to complete mission: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
