package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbDSN = "host=localhost port=54321 user=bench database=bench password=bench"

func main() {
	ctx := context.Background()
	pool := connect(ctx, dbDSN)
	defer pool.Close()

	initDB(pool)

	go listen(context.Background(), pool)

	eventBody := `{"id": 1, "name": "test"}`
	for i := 0; ; i++ {
		_, err := pool.Exec(ctx, "INSERT INTO events (body) VALUES ($1)", eventBody)
		if err != nil {
			log.Println("Error sending notification:", err)
			os.Exit(1)
		}
	}
}

func connect(ctx context.Context, dsn string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println("failed to parse DB DSN:", err)
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Println("failed to connect to DB:", err)
		os.Exit(1)
	}
	return pool
}

func initDB(pool *pgxpool.Pool) {
	fCont, err := os.ReadFile("init.sql")
	if err != nil {
		panic(err)
	}
	sql := string(fCont)
	_, err = pool.Exec(context.Background(), sql)
	if err != nil {
		panic(err)
	}
}

func listen(ctx context.Context, pool *pgxpool.Pool) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println("Error acquiring connection:", err)
		os.Exit(1)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "listen events")
	if err != nil {
		log.Println("Error listening to chat channel:", err)
		os.Exit(1)
	}

	t := time.Now()
	var i int
	for {
		_, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			log.Println("Error waiting for notification:", err)
			os.Exit(1)
		}
		i++
		if time.Since(t) > time.Second {
			fmt.Println(i)
			t = time.Now()
			i = 0
		}
	}
}
