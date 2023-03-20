package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/artemmarkaryan/exlex-backend/internal/migrations"
	"github.com/artemmarkaryan/exlex-backend/internal/server"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	{
		sigc := make(chan os.Signal, 1)
		signal.Notify(
			sigc,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		go func() {
			<-sigc
			cancel()
		}()
	}

	ctx = connectDatabase(ctx)
	if err := server.Serve(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("finished")
}

func connectDatabase(ctx context.Context) context.Context {
	dsn := os.Getenv("PSQL_DSN")
	if dsn == "" {
		log.Fatalln("PSQL_DSN is not set")
	}

	ctx = database.Propagate(ctx, database.Connect(ctx, dsn))

	db := database.C(ctx)
	if err := db.Ping(); err != nil {
		log.Fatalf("database ping: %v\n", err)
	}

	if err := goose.Up(db.DB, "/"); err != nil {
		log.Fatalln(err)
	}

	return ctx
}
