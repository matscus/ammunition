package main

import (
	"ammunition/cache"
	"ammunition/config"
	"ammunition/database"
	"ammunition/docs"
	"ammunition/handlers"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	hostIp, configPath, actuatorPath, pemPath, keyPath, proto, listenport, dbuser, dbpassword, dbhost, dbname, logLevel string
	dbport                                                                                                              int
	wait, writeTimeout, readTimeout, idleTimeout                                                                        time.Duration
	debug, logger                                                                                                       bool
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	flag.StringVar(&pemPath, "pempath", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "keypath", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", "9443", "port to Listen")
	flag.StringVar(&proto, "proto", "http", "http or https")
	flag.StringVar(&hostIp, "host-ip", "", "set current ip adress")
	flag.StringVar(&configPath, "cache-config", "config.yaml", "path from cache config file")
	flag.StringVar(&actuatorPath, "actuator-config", "actuator.yaml", "path from actuator config file")
	flag.StringVar(&dbuser, "dbuser", "postgres", "db user")
	flag.StringVar(&dbpassword, "dbpassword", `postgres`, "db user password")
	flag.StringVar(&dbhost, "dbhost", "localhost", "db host")
	flag.StringVar(&logLevel, "loglevel", "INFO", "log level, default INFO")
	flag.StringVar(&docs.SwaggerConfPath, "swagger", "swagger.yaml", "path to swagger config file")
	flag.IntVar(&dbport, "dbport", 5432, "db port")
	flag.StringVar(&dbname, "dbname", "ammunition", "db name")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*60, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*60, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.Parse()
	err := config.ReadCacheConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = config.ReadActuatorConfig(actuatorPath)
	if err != nil {
		log.Fatal(err)
	}
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	var router *gin.Engine
	if logger {
		router = gin.Default()

	} else {
		router = gin.New()
		router.Use(gin.Recovery())
	}
	router.Use(handlers.Middleware())
	ppof := router.Group("/pprof")
	ppof.GET("/", gin.WrapF(pprof.Index))
	ppof.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	ppof.GET("/profile", gin.WrapF(pprof.Profile))
	ppof.POST("/symbol", gin.WrapF(pprof.Symbol))
	ppof.GET("/symbol", gin.WrapF(pprof.Symbol))
	ppof.GET("/trace", gin.WrapF(pprof.Trace))
	ppof.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	ppof.GET("/block", gin.WrapH(pprof.Handler("block")))
	ppof.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	ppof.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	ppof.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	ppof.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))

	v2 := router.Group("/api/v2")
	{
		actuator := v2.Group("actuator")
		{
			actuator.GET("/metrics", handlers.PrometheusHandler())
			actuator.GET("/info", handlers.Info)
			actuator.GET("/health", handlers.Health)
		}
		cache := v2.Group("cache")
		{
			cache.Any("/", handlers.DefaultHandle)
			cache.Any("/persisted", handlers.PersistHandle)
			temporary := cache.Group("temporary")
			{
				temporary.Any("/", handlers.TemporaryHandle)
				temporary.Any("/init", handlers.TemporaryInitHandle)
			}
		}
		docs.SwaggerInfo.BasePath = "/api/v2/"
		v2.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	go func() {
		for {
			err := database.InitDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname))
			if err != nil {
				log.Error(err)
			} else {
				err = cache.InitAllPersistedPools()
				if err != nil {
					log.Error(err)
				} else {
					log.Info("Init all persist pool completed")
				}
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
	cache.InitDefault()
	log.Info("Init Default pool completed")
	if hostIp == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Error("Get interface adress error: ", err.Error())
			os.Exit(1)
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					hostIp = ipnet.IP.String()
				}
			}
		}
	}
	srv := &http.Server{
		Addr:         hostIp + ":" + listenport,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      router,
	}
	go func() {
		switch proto {
		case "https":
			log.Info("Server is run, proto: https, address: %s ", srv.Addr)
			if err := srv.ListenAndServeTLS(pemPath, keyPath); err != nil {
				log.Println(err)
			}
		case "http":
			log.Info("Server is run, proto: http, address: %s ", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}
	}()
	logo()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Info("server shutting down")
	os.Exit(0)
}

func logo() {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("\r%s%s\n", "\033[;31m", `   
   ▄████████   ▄▄▄▄███▄▄▄▄     ▄▄▄▄███▄▄▄▄   ███    █▄  ███▄▄▄▄    ▄█      ███      ▄█   ▄██████▄  ███▄▄▄▄   
  ███    ███ ▄██▀▀▀███▀▀▀██▄ ▄██▀▀▀███▀▀▀██▄ ███    ███ ███▀▀▀██▄ ███  ▀█████████▄ ███  ███    ███ ███▀▀▀██▄ 
  ███    ███ ███   ███   ███ ███   ███   ███ ███    ███ ███   ███ ███▌    ▀███▀▀██ ███▌ ███    ███ ███   ███ 
  ███    ███ ███   ███   ███ ███   ███   ███ ███    ███ ███   ███ ███▌     ███   ▀ ███▌ ███    ███ ███   ███ 
▀███████████ ███   ███   ███ ███   ███   ███ ███    ███ ███   ███ ███▌     ███     ███▌ ███    ███ ███   ███ 
  ███    ███ ███   ███   ███ ███   ███   ███ ███    ███ ███   ███ ███      ███     ███  ███    ███ ███   ███ 
  ███    ███ ███   ███   ███ ███   ███   ███ ███    ███ ███   ███ ███      ███     ███  ███    ███ ███   ███ 
  ███    █▀   ▀█   ███   █▀   ▀█   ███   █▀  ████████▀   ▀█   █▀  █▀      ▄████▀   █▀    ▀██████▀   ▀█   █▀  
  `)
}
