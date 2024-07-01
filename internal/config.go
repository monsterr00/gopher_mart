package internal

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	Server     Server
	Db         Db
	AccrualSys AccrualSys
	Resty      Resty
}

type Server struct {
	Host string
}

type Db struct {
	Address        string
	MigrationsPath string
	Conn           *sql.DB
}

type AccrualSys struct {
	Host string
}

type Resty struct {
	Client *resty.Client
}

var Flags struct {
	ServerHost        string
	DBHost            string
	AccrualSystemHost string
}

func LoadConfig() (Config, error) {
	var config Config
	var err error

	config.Server, err = loadServer()
	if err != nil {
		/// поменять обработку ошибок
		return config, fmt.Errorf("could not load server config: %w", err)
	}

	config.Db, err = loadDB()
	if err != nil {
		/// поменять обработку ошибок
		return config, fmt.Errorf("could not load db config: %w", err)
	}

	config.AccrualSys, err = loadAccuralSys()
	if err != nil {
		/// поменять обработку ошибок
		return config, fmt.Errorf("could not load accrual system config: %w", err)
	}

	config.Resty, err = loadResty()
	if err != nil {
		/// поменять обработку ошибок
		return config, fmt.Errorf("could not resty config: %w", err)
	}

	return config, nil
}

func loadServer() (Server, error) {
	envAddress, isSet := os.LookupEnv("RUN_ADDRESS")
	if isSet && envAddress != "" {
		Flags.ServerHost = envAddress
	}

	return Server{
		Host: Flags.ServerHost,
	}, nil
}

func loadDB() (Db, error) {
	envAddress, isSet := os.LookupEnv("DATABASE_URI")
	if isSet && envAddress != "" {
		Flags.DBHost = envAddress
	}

	return Db{
		Address:        Flags.DBHost,
		MigrationsPath: "db/migrations",
		Conn:           nil,
	}, nil
}

func loadAccuralSys() (AccrualSys, error) {
	envAddress, isSet := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if isSet && envAddress != "" {
		Flags.AccrualSystemHost = envAddress
	}

	return AccrualSys{
		Host: Flags.AccrualSystemHost,
	}, nil
}

func loadResty() (Resty, error) {
	restyClient := resty.New()

	restyClient.
		SetRetryCount(3).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)
	return Resty{
		Client: restyClient,
	}, nil
}
