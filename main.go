package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"kickof/config"
	"kickof/controllers"
	"kickof/database"
	"kickof/models"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	if !database.Init() {
		log.Printf("Connected to MongoDB URI: Failure")
		return
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"POST", "GET", "PATCH", "OPTIONS", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")
	{
		api.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Data: "KickOf Api v" + os.Getenv("VERSION"),
			})
			return
		})

		api.POST("/register", controllers.SignUp)
		api.POST("/login", controllers.SignIn)
		api.GET("/refresh-token", controllers.RefreshToken)
		api.POST("/activate", controllers.Activate)

		protected := api.Group("/", config.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)

			protected.GET("/project", controllers.GetProjects)
			protected.POST("/project", controllers.CreateProject)
			protected.GET("/project/:id", controllers.GetProjectByIdOrCode)
			protected.PATCH("/project/:id", controllers.UpdateProject)
			protected.DELETE("/project/:id", controllers.DeleteProject)

			protected.GET("/state", controllers.GetStates)
			protected.POST("/state", controllers.CreateState)
			protected.GET("/state/:id", controllers.GetStateById)
			protected.PATCH("/state/:id", controllers.UpdateState)
			protected.DELETE("/state/:id", controllers.DeleteState)

			protected.GET("/task", controllers.GetTasks)
			protected.POST("/task", controllers.CreateTask)
			protected.GET("/task/:id", controllers.GetTaskById)
			protected.PATCH("/task/:id", controllers.UpdateTask)
			protected.DELETE("/task/:id", controllers.DeleteTask)

			protected.GET("/task-label", controllers.GetTaskLabels)
			protected.POST("/task-label", controllers.CreateTaskLabel)
			protected.GET("/task-label/:id", controllers.GetTaskLabelById)
			protected.PATCH("/task-label/:id", controllers.UpdateTaskLabel)
			protected.DELETE("/task-label/:id", controllers.DeleteTaskLabel)

			protected.GET("/workspace", controllers.GetWorkspaces)
			protected.POST("/workspace", controllers.CreateWorkspace)
			protected.GET("/workspace/members/:workspaceId", controllers.GetWorkspaceMembers)
			protected.GET("/workspace/:id", controllers.GetWorkspaceById)
			protected.PATCH("/workspace/:id", controllers.UpdateWorkspace)
			protected.DELETE("/workspace/:id", controllers.DeleteWorkspace)
		}
	}

	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	err = router.Run(":" + port)
	if err != nil {
		return
	}
}
