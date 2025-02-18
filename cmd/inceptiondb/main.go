package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fulldump/box"
	"github.com/fulldump/goconfig"

	"github.com/fulldump/inceptiondb/api"
	"github.com/fulldump/inceptiondb/configuration"
	"github.com/fulldump/inceptiondb/database"
	"github.com/fulldump/inceptiondb/service"
)

var VERSION = "dev"

var banner = `
 _____                     _   _            ____________ 
|_   _|                   | | (_)           |  _  \ ___ \
  | | _ __   ___ ___ _ __ | |_ _  ___  _ __ | | | | |_/ /
  | || '_ \ / __/ _ \ '_ \| __| |/ _ \| '_ \| | | | ___ \
 _| || | | | (_|  __/ |_) | |_| | (_) | | | | |/ /| |_/ /
 \___/_| |_|\___\___| .__/ \__|_|\___/|_| |_|___/ \____/ 
                    | |                                  
                    |_|          version ` + VERSION + `
`

func main() {

	c := configuration.Default()
	goconfig.Read(&c)

	if c.Version {
		fmt.Println("Version:", VERSION)
		return
	}

	if c.ShowBanner {
		fmt.Println(banner)
	}

	if c.ShowConfig {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "    ")
		e.Encode(c)
	}

	db := database.NewDatabase(&database.Config{
		Dir: c.Dir,
	})

	b := api.Build(service.NewService(db), c.Statics, VERSION)
	accessLogger := log.New(os.Stdout, "ACCESS: ", log.Lshortfile)
	b.WithInterceptors(
		api.AccessLog(accessLogger),
		api.InterceptorUnavailable(db),
		api.RecoverFromPanic,
		api.PrettyErrorInterceptor,
	)

	s := &http.Server{
		Addr:    c.HttpAddr,
		Handler: box.Box2Http(b),
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			sig := <-signalChan
			fmt.Println("Signal received", sig.String())
			db.Stop()
			s.Shutdown(context.Background())
		}
	}()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := db.Start()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("listening on", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	wg.Wait()
}
