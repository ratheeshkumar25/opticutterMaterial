package di

import (
	"log"

	"github.com/ratheeshkumar25/opt_cut_material_service/config"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/db"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/handlers"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/repo"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/server"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/service"
)

func Init() {
	cnfg := config.LoadConfig()

	db := db.ConnectDB(cnfg)

	materialRepo := repo.NewMaterialRepository(db)
	materialService := service.NewMaterialService(materialRepo)
	materialHanldr := handlers.NewMaterialHandler(materialService)

	err := server.NewGrpcMaterialServer(cnfg.GrpcPort, materialHanldr)
	if err != nil {
		log.Fatalf("failed to start gRPC server %v", err.Error())
	}
}
