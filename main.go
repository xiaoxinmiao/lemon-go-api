package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"

	"sample-go-api/factory"
)

var (
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func main() {
	appEnv := flag.String("app-env", os.Getenv("APP_ENV"), "app env")
	fruitConnEnv := flag.String("FRUIT_CONN", os.Getenv("FRUIT_CONN"), "FRUIT_CONN")
	flag.Parse()

	var c Config
	if err := configutil.Read(*appEnv, &c); err != nil {
		panic(err)
	}

	fmt.Println(c)
	db, err := initDB("mysql", *fruitConnEnv)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.GET("/swagger", func(c echo.Context) error {
		return c.File("./swagger.yml")
	})
	e.Static("/swagger-yml", "./swagger-yml")
	e.Static("/docs", "./swagger-ui")
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.RequestID())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDB(c.Service, db, echomiddleware.KafkaConfig(c.Logger.Kafka)))
	e.Use(echomiddleware.BehaviorLogger(c.Service, echomiddleware.KafkaConfig(c.BehaviorLog.Kafka)))

	e.Validator = &Validator{}

	e.Debug = c.Debug

	configMap := map[string]interface{}{
		"key": "123",
	}
	setContextValueMiddleware := setContextValue(&configMap)
	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return setContextValueMiddleware(handlerFunc)(c)
	}
	e.Use(setContextValueMiddleware)
	if err := e.Start(":" + c.HttpPort); err != nil {
		log.Println(err)
	}
}

func setContextValue(configMap *map[string]interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			reqContext := context.WithValue(req.Context(), factory.ContextConfigName, configMap)
			c.SetRequest(req.WithContext(reqContext))
			return next(c)
		}
	}
}

func initDB(driver, connection string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Config struct {
	Logger struct {
		Kafka echomiddleware.KafkaConfig
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	Trace struct {
		Zipkin echomiddleware.ZipkinConfig
	}

	Debug    bool
	Service  string
	HttpPort string
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}