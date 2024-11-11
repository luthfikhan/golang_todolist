package main

import (
	"context"
	"net/http"
	"todolist/app"
	"todolist/controller"
	"todolist/helper"
	"todolist/middleware"
	"todolist/repository"
	"todolist/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	err := godotenv.Load()
	if err != nil {
		helper.Log.Error(nil, err)
	}
}

func main() {
	router := gin.New()
	context := context.Background()

	router.Use(middleware.Log, middleware.Recover())
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})

	db := app.GetOracleDb()
	defer db.Close()
	redisClient := app.InitiateRedis()
	defer redisClient.Close()

	cache := &helper.Cache{
		Ctx: context,
		Rdb: redisClient,
	}

	tasksRepository := repository.NewTasksRepository(db)
	tasksService := service.NewTasksService(tasksRepository, cache)
	tasksController := controller.NewTasksController(tasksService)

	tasksRoute := router.Group("/tasks")
	tasksRoute.Use(middleware.Auth)

	tasksRoute.POST("", tasksController.InsertTask)
	tasksRoute.GET("", tasksController.GetAllTasks)
	tasksRoute.GET("/:id", tasksController.GetTask)
	tasksRoute.PUT("/:id", tasksController.UpdateTask)
	tasksRoute.DELETE("/:id", tasksController.DeleteTask)

	router.GET("/token", func(ctx *gin.Context) {
		token, err := helper.CreateToken()
		helper.PanicIfError(err)

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	})

	router.Run()
}
