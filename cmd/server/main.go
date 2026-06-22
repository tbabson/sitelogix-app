package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/internal/attendance"
	"github.com/sitelogix/backend/internal/auth"
	"github.com/sitelogix/backend/internal/dailylog"
	"github.com/sitelogix/backend/internal/dashboard"
	"github.com/sitelogix/backend/internal/inventory"
	"github.com/sitelogix/backend/internal/issues"
	"github.com/sitelogix/backend/internal/materials"
	"github.com/sitelogix/backend/internal/media"
	"github.com/sitelogix/backend/internal/notification"
	"github.com/sitelogix/backend/internal/project"
	"github.com/sitelogix/backend/internal/reports"
	"github.com/sitelogix/backend/internal/requisition"
	"github.com/sitelogix/backend/internal/search"
	"github.com/sitelogix/backend/internal/user"
	"github.com/sitelogix/backend/pkg/config"
	"github.com/sitelogix/backend/pkg/database"
	"github.com/sitelogix/backend/pkg/fcm"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/storage"
	"github.com/sitelogix/backend/pkg/weather"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()
	log.Println("connected to postgres")

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// External clients
	gh := storage.NewGitHubStorage(cfg.GitHub)
	weatherClient := weather.NewClient(cfg.Weather.APIKey, cfg.Weather.APIURL)
	fcmClient := fcm.NewClient(cfg.FCM.ServerKey)

	// Repositories
	userRepo := user.NewRepository(db)
	authRepo := auth.NewRepository(db)

	// Services
	authSvc := auth.NewService(authRepo, userRepo, cfg.JWT)
	userSvc := user.NewService(db)
	projectSvc := project.NewService(db)
	invSvc := inventory.NewService(db)
	materialsSvc := materials.NewService(db)
	reqSvc := requisition.NewService(db, invSvc)
	logSvc := dailylog.NewService(db, weatherClient, invSvc)
	attendanceSvc := attendance.NewService(db)
	issueSvc := issues.NewService(db)
	mediaSvc := media.NewService(db, gh)
	notifSvc := notification.NewService(db, fcmClient)
	dashboardSvc := dashboard.NewService(db)
	reportsSvc := reports.NewService(db)
	searchSvc := search.NewService(db)

	// Handlers
	authHandler := auth.NewHandler(authSvc)
	userHandler := user.NewHandler(userSvc)
	projectHandler := project.NewHandler(projectSvc)
	invHandler := inventory.NewHandler(invSvc)
	materialsHandler := materials.NewHandler(materialsSvc)
	reqHandler := requisition.NewHandler(reqSvc)
	logHandler := dailylog.NewHandler(logSvc, notifSvc)
	attendanceHandler := attendance.NewHandler(attendanceSvc)
	issueHandler := issues.NewHandler(issueSvc, notifSvc)
	mediaHandler := media.NewHandler(mediaSvc)
	notifHandler := notification.NewHandler(notifSvc)
	dashboardHandler := dashboard.NewHandler(dashboardSvc)
	reportsHandler := reports.NewHandler(reportsSvc)
	searchHandler := search.NewHandler(searchSvc)

	// Router
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "app": "sitelogix"})
	})

	api := r.Group("/api/v1")
	auth.RegisterRoutes(api, authHandler)
	user.RegisterRoutes(api, userHandler, cfg.JWT.Secret)
	project.RegisterRoutes(api, projectHandler, cfg.JWT.Secret)
	materials.RegisterRoutes(api, materialsHandler, cfg.JWT.Secret)
	inventory.RegisterRoutes(api, invHandler, cfg.JWT.Secret)
	requisition.RegisterRoutes(api, reqHandler, cfg.JWT.Secret)
	dailylog.RegisterRoutes(api, logHandler, cfg.JWT.Secret)
	attendance.RegisterRoutes(api, attendanceHandler, cfg.JWT.Secret)
	issues.RegisterRoutes(api, issueHandler, cfg.JWT.Secret)
	media.RegisterRoutes(api, mediaHandler, cfg.JWT.Secret)
	notification.RegisterRoutes(api, notifHandler, cfg.JWT.Secret)
	dashboard.RegisterRoutes(api, dashboardHandler, cfg.JWT.Secret)
	reports.RegisterRoutes(api, reportsHandler, cfg.JWT.Secret)
	search.RegisterRoutes(api, searchHandler, cfg.JWT.Secret)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("SiteLogix API running on %s [%s]", addr, cfg.App.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
