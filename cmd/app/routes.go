package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"spy-cats/internal/database"
	catRepository "spy-cats/internal/repository/cat"
	missionRepository "spy-cats/internal/repository/missions"
	catService "spy-cats/internal/service/cats"
	missionService "spy-cats/internal/service/missions"
	cat_handler "spy-cats/internal/transport/rest/cat"
	mission_handler "spy-cats/internal/transport/rest/missions"
)

func SetupRoutes(ginEngine *gin.Engine) {
	log := SetupLogger()
	catRoutes(ginEngine, log)
	missionRoutes(ginEngine, log)
	targetRoutes(ginEngine, log)
}

func catRoutes(ginEngine *gin.Engine, log *slog.Logger) {
	r := ginEngine.Group("/cats/")
	db := database.NewPGClient()

	repository := catRepository.NewCats(db.DB)
	service := catService.NewCats(repository)
	handler := cat_handler.NewHandler(service, log)

	r.GET("", handler.GetAll)
	r.GET(":id", handler.GetById)
	r.POST("create", handler.Create)
	r.DELETE(":id/delete", handler.Delete)
	r.PATCH(":id/update", handler.Update)
}

func missionRoutes(ginEngine *gin.Engine, log *slog.Logger) {
	r := ginEngine.Group("/missions/")
	db := database.NewPGClient()

	missionRepo := missionRepository.NewMissions(db.DB)
	targetRepo := missionRepository.NewTargets(db.DB)
	service := missionService.NewMissions(missionRepo, targetRepo)
	handler := mission_handler.NewHandler(service, log)

	r.GET("", handler.GetAll)
	r.GET(":id", handler.GetById)
	r.POST("create", handler.Create)
	r.DELETE(":id/delete", handler.Delete)
	r.PATCH(":id/assign-cat/:cat_id", handler.AssignCat)
	r.PATCH(":id/complete", handler.CompleteMission)
}

func targetRoutes(ginEngine *gin.Engine, log *slog.Logger) {
	r := ginEngine.Group("/targets/")
	db := database.NewPGClient()

	missionRepo := missionRepository.NewMissions(db.DB)
	targetRepo := missionRepository.NewTargets(db.DB)
	service := missionService.NewTargets(missionRepo, targetRepo)
	handler := mission_handler.NewTargetHandler(service, log)

	r.POST("create", handler.Create)
	r.PATCH(":id/add/:mission_id", handler.AddTargetToMission)
	r.PATCH(":id/delete", handler.DeleteTargetFromMission)
	r.PATCH(":id/complete", handler.CompleteTarget)
	r.PATCH(":id/update-notes", handler.UpdateNotes)
}
