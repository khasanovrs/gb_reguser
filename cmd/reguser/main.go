package main

import (
	"context"
	"github.com/khasanovrs/gb_reguser/api/handler"
	"github.com/khasanovrs/gb_reguser/api/routeroapi"
	"github.com/khasanovrs/gb_reguser/api/server"
	"github.com/khasanovrs/gb_reguser/app/repos/user"
	"github.com/khasanovrs/gb_reguser/app/starter"
	"github.com/khasanovrs/gb_reguser/db/mem/usermemstore"
	"github.com/khasanovrs/gb_reguser/db/sql/pgstore"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	// output current time zone
	tnow := time.Now()
	tz, _ := tnow.Zone()
	log.Printf("Local time zone %s. Service started at %s", tz,
		tnow.Format("2006-01-02T15:04:05.000 MST"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	var ust user.UserStore
	stu := os.Getenv("REGUSER_STORE")

	switch stu {
	case "mem":
		ust = usermemstore.NewUsers()
	case "pg":
		dsn := os.Getenv("DATABASE_URL")
		pgst, err := pgstore.NewUsers(dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer pgst.Close()
		ust = pgst
	default:
		log.Fatal("unknown REGUSER_STORE = ", stu)
	}

	a := starter.NewApp(ust)
	us := user.NewUsers(ust)
	h := handler.NewHandlers(us)

	rh := routeroapi.NewRouterOpenAPI(h)

	srv := server.NewServer(":"+os.Getenv("PORT"), rh)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
