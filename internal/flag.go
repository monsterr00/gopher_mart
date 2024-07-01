package internal

import "flag"

func init() {
	flag.StringVar(&Flags.ServerHost, "a", "localhost:8080", "server host")
	flag.StringVar(&Flags.DBHost, "d", "host=localhost user=postgres password=postgres1 dbname=metrics sslmode=disable", "db address")
	flag.StringVar(&Flags.AccrualSystemHost, "r", "localhost:8081", "accrual system host")
}
