package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"context"

	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/matscus/ammunition/cache"
	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/handlers"
	"github.com/matscus/ammunition/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	pemPath, keyPath, proto, listenport, host, dbuser, dbpassword, dbhost, dbname, logLevel string
	dbport                                                                                  int
	wait, writeTimeout, readTimeout, idleTimeout                                            time.Duration
)

func main() {

	flag.StringVar(&pemPath, "pempath", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "keypath", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", "10000", "port to Listen")
	flag.StringVar(&proto, "proto", "http", "http or https")
	flag.StringVar(&dbuser, "user", "postgres", "db user")
	flag.StringVar(&dbpassword, "password", `postgres`, "db user password")
	flag.StringVar(&dbhost, "host", "localhost", "db host")
	flag.StringVar(&logLevel, "loglevel", "INFO", "log level, default INFO")
	flag.IntVar(&dbport, "dbport", 5432, "db port")
	flag.StringVar(&dbname, "dbname", "ammunition", "db name")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*15, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*15, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.Parse()
	log.Info("Parse flag completed")
	setLogLevel(logLevel)
	log.Info("Set log level completed")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/cache/persisted", middleware.Middleware(handlers.PersistedDatapoolHandler)).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/cache/cookies", middleware.Middleware(handlers.CookiesHandler)).Methods(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/cache/kv", middleware.Middleware(handlers.KVHahdler)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	http.Handle("/", r)
	log.Info("Register handlers and route completed")
	go func() {
		for {
			err := database.InitDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname))
			if err != nil {
				log.Error(err)
			} else {
				err = cache.InitAllPersistedPools()
				if err != nil {
					log.Error(err)
				}
				break
			}
			time.Sleep(10 * time.Second)
		}
		log.Info("Start init persisted pools")
	}()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("Get interface adress error: ", err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
			}
		}
	}
	log.Info("Get IPv4 addr completed")
	srv := &http.Server{
		Addr:         host + ":" + listenport,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      r,
	}
	log.Info("Set server params completed")
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
	srv.Shutdown(ctx)
	log.Info("server shutting down")
	os.Exit(0)
}

func setLogLevel(level string) {
	level = strings.ToUpper(level)
	switch level {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	}
}

func logo() {
	time.Sleep(50 * time.Millisecond)
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
