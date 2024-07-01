package internal

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	app "github.com/monsterr00/gopher_mart/internal/application"
	"github.com/monsterr00/gopher_mart/internal/helpers"
	ordersInfra "github.com/monsterr00/gopher_mart/internal/infrastructure/orders"
	usersInfra "github.com/monsterr00/gopher_mart/internal/infrastructure/users"
	withdrawalsInfra "github.com/monsterr00/gopher_mart/internal/infrastructure/withdrawals"
	oCrServ "github.com/monsterr00/gopher_mart/internal/service/orders"
	uCrServ "github.com/monsterr00/gopher_mart/internal/service/users"
	wCrServ "github.com/monsterr00/gopher_mart/internal/service/withdrawals"
)

func Run() error {
	flag.Parse()

	cfg, err := LoadConfig()
	if err != nil {
		// Поменять обработку ошибок
		//slog.Error("Could not load config", "err", err)
		return err
	}

	cfg.Db.Conn = startDB(cfg)
	startServer(cfg)

	return nil
}

func startServer(cfg Config) {
	userRepo := usersInfra.NewUserPostgresRepo(cfg.Db.Conn)
	orderRepo := ordersInfra.NewOrderPostgresRepo(cfg.Db.Conn)
	withdrawalRepo := withdrawalsInfra.NewWithdrawalPostgresRepo(cfg.Db.Conn)

	uCrRepo := uCrServ.NewUserCreationService(userRepo)
	oCrRepo := oCrServ.NewOrderCreationService(orderRepo, userRepo, cfg.AccrualSys.Host, cfg.Resty.Client)
	wCrRepo := wCrServ.NewWithdrawalCreationService(orderRepo, withdrawalRepo)

	server := app.NewServer(uCrRepo, oCrRepo, wCrRepo)

	err := http.ListenAndServe(cfg.Server.Host, server)
	if err != nil {
		// Поменять обработку ошибок
		log.Fatal(err)
	}
}

func startDB(cfg Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.Db.Address)
	if err != nil {
		// Поменять обработку ошибок
		log.Fatal(err)
	}

	filePath := helpers.AbsolutePath("file:///", cfg.Db.MigrationsPath)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		// Поменять обработку ошибок
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		filePath,
		"postgres", driver)

	if err != nil {
		// Поменять обработку ошибок
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	err = db.Ping()
	if err != nil {
		// Поменять обработку ошибок
		log.Fatal(err)
	}

	return db
}

/*
func (storl *store) Close() error {
	return storl.conn.Close()
}
*/
