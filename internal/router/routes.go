package router

import (
	"net/http"

	log "github.com/go-kit/log"

	"github.com/Quaqmre/go-case/internal/handler/fetch"
	"github.com/Quaqmre/go-case/internal/handler/inmemory"
	"github.com/Quaqmre/go-case/internal/infrastructure/persistent"
	"github.com/Quaqmre/go-case/internal/middleware"
)

func RegisterRoutes(storage inmemory.MemStore, mongoStorage persistent.Query, logger log.Logger) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/fetch", middleware.LoggerMiddleWare(logger)(fetch.NewHandler(mongoStorage)))

	router.Handle("/in-memory", middleware.LoggerMiddleWare(logger)(inmemory.NewHandler(storage)))

	return router
}
