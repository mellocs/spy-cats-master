package cat

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"spy-cats/internal/models/cats"
	catService "spy-cats/internal/service/cats"
	"strconv"
)

type Handler struct {
	catService catService.Cats
	logger     *slog.Logger
}

func NewHandler(catService catService.Cats, logger *slog.Logger) *Handler {
	return &Handler{
		catService: catService,
		logger:     logger,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	allCats, err := h.catService.GetAll()
	if err != nil {
		h.logger.Error("failed to get all cats: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all cats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": allCats})
}

func (h *Handler) Create(c *gin.Context) {
	var cat cats.Cat
	err := c.BindJSON(&cat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat JSON"})
		return
	}

	if cat.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat name"})
		return
	}
	if cat.YearsOfExperience < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "years of experience cannot be negative"})
		return
	}

	if err = h.catService.Create(cat); err != nil {
		h.logger.Error("failed to create cat: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *Handler) Update(c *gin.Context) {
	catID, err := strconv.Atoi(c.Param("id"))
	if catID < 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat ID"})
		return
	}

	var input struct {
		Salary uint `json:"salary"`
	}

	if err = c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if input.Salary <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "salary must be greater than zero"})
		return
	}

	if err = h.catService.Update(catID, input.Salary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h *Handler) Delete(c *gin.Context) {
	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}

	err = h.catService.Delete(catID)
	if err != nil {
		h.logger.Error("failed to delete cat: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}

func (h *Handler) GetById(c *gin.Context) {
	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}
	cat, err := h.catService.GetByID(catID)
	if errors.Is(err, cats.ErrCatNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		h.logger.Error("failed to get cat: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": cat})
}
