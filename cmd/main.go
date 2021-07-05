package main

import (
	"context"
	"log"
	"os"

	"github.com/cristoper/gsheet/gdrive"
)

// global FilesService object
var svc = func() *gdrive.Service {
	svc, err := gdrive.NewServiceWithCtx(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return svc
}()

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
