package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"context"

	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/handlers"
	"github.com/matscus/ammunition/middleware"
	"github.com/matscus/ammunition/pool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	pemPath, keyPath, proto, listenport, host, dbuser, dbpassword, dbhost, dbname string
	dbport                                                                        int
	wait, writeTimeout, readTimeout, idleTimeout                                  time.Duration
)

func main() {

	flag.StringVar(&pemPath, "pempath", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "keypath", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", "10000", "port to Listen")
	flag.StringVar(&proto, "proto", "http", "http or https")
	flag.StringVar(&dbuser, "user", "postgres", "db user")
	flag.StringVar(&dbpassword, "password", `postgres`, "db user password")
	flag.StringVar(&dbhost, "host", "postgres", "db host")
	flag.IntVar(&dbport, "dbport", 5432, "db port")
	flag.StringVar(&dbname, "dbname", "postgres", "db name")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*15, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*15, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/persisted/manage", middleware.Middleware(handlers.PersistedManageHandler)).Methods(http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/persisted", middleware.Middleware(handlers.PersistedGetHandler)).Methods(http.MethodGet, http.MethodOptions).Queries("name", "{name}", "project", "{project}")
	//r.HandleFunc("/api/v1/datapool/temporary", middleware.Middleware(handlers.GetValue)).Methods(http.MethodPost, http.MethodOptions)
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

	go func() {
		for {
			err := database.InitDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname))
			if err == nil {
				err = pool.InitAllPersistedPools()
				if err != nil {
					log.Println(err)
				}
				break
			}
		}
	}()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("Get interface adres error: ", err.Error())
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
			}
		}
	}
	srv := &http.Server{
		Addr:         host + ":" + listenport,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      r,
	}

	go func() {
		switch proto {
		case "https":
			log.Printf("Server is run, proto: https, address: %s ", srv.Addr)
			if err := srv.ListenAndServeTLS(pemPath, keyPath); err != nil {
				log.Println(err)
			}
		case "http":
			log.Printf("Server is run, proto: http, address: %s ", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}

	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("server shutting down")
	os.Exit(0)
}
