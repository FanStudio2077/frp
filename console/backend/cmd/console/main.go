package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FanStudio2077/frp-console/internal/api"
	"github.com/FanStudio2077/frp-console/internal/config"
)

func main() {
	cfg := config.Load()
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: api.NewRouter()}
	go func() {
		log.Printf("frp-console listening on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("console server error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
