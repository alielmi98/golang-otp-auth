package main

import (
	"fmt"
	"log"

	"github.com/alielmi98/golang-otp-auth/docs"
	_ "github.com/alielmi98/golang-otp-auth/docs"
	"github.com/alielmi98/golang-otp-auth/internal/middlewares"
	"github.com/alielmi98/golang-otp-auth/internal/user/api/handler"
	usersRouter "github.com/alielmi98/golang-otp-auth/internal/user/api/router"
	"github.com/alielmi98/golang-otp-auth/internal/user/api/validation"
	"github.com/alielmi98/golang-otp-auth/migrations"
	"github.com/alielmi98/golang-otp-auth/pkg/cache"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/alielmi98/golang-otp-auth/pkg/constants"
	"github.com/alielmi98/golang-otp-auth/pkg/db"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {

	cfg := config.GetConfig()

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Redis, constants.Startup, err.Error())
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Postgres, constants.Startup, err.Error())
	}

	migrations.Up1()
	InitServer(cfg)

}
func InitServer(cfg *config.Config) {
	r := gin.New()
	RegisterValidators()

	userHandler := handler.NewUserHandler(cfg)
	r.Use(middlewares.Cors(cfg))
	RegisterRoutes(r, cfg, userHandler)
	RegisterSwagger(r, cfg)
	log.Printf("Caller:%s Level:%s Msg:%s", constants.General, constants.Startup, "Started")
	r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))

}

func RegisterRoutes(r *gin.Engine, cfg *config.Config, userHandler *handler.UsersHandler) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		//Auth
		users := v1.Group("/users")
		usersRouter.Users(users, cfg, userHandler)

	}
}

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			log.Printf("Caller:%s Level:%s Msg:%s", constants.Validation, constants.Startup, err.Error())
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "This is a sample auth by otp for golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.InternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
