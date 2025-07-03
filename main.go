package main

import (
	"fmt"
	"log"
	"os"
	"supplier-go-service/db"
	"supplier-go-service/internal/supplier"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	godotenv.Load()
	natsServer := getEnv("NATS_SERVERS", nats.DefaultURL)
	fmt.Println("\033[33m🔌 Conectando a NATS...\033[0m")
	nc, err := nats.Connect(natsServer)
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a NATS: %v", err)
	}
	fmt.Println("\033[32m✅ Conexión a NATS exitosa.\033[0m")

	user := getEnv("DB_USERNAME", "")
	pass := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "")
	port := getEnv("DB_PORT", "")
	dbName := getEnv("DB_NAME", "")
	dataBase := db.InitDB(user, pass, host, port, dbName)
	repo := supplier.NewRepository(dataBase)
	service := supplier.NewService(repo)
	natsH := supplier.NewHandler(service)
	natsH.Subscribe(nc)

	select {} // Prevent main from exiting
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
