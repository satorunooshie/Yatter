package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/satorunooshie/Yatter/app/app"
	"github.com/satorunooshie/Yatter/app/config"
	"github.com/satorunooshie/Yatter/app/handler"
)

func main() {
	log.Fatalf("%+v", serve(context.Background()))
}

func serve(ctx context.Context) error {
	app, err := app.NewApp()
	if err != nil {
		return err
	}
	addr := ":" + strconv.Itoa(config.Port())
	log.Printf("Serve on http://%s", addr)

	return http.ListenAndServe(addr, handler.NewRouter(app))
}
