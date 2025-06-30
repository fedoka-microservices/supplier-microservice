package main

import (
	"fmt"
	"log"
	"supplier-go-service/db"
	"supplier-go-service/internal/supplier"

	"github.com/nats-io/nats.go"
)

func main() {
	db := db.InitDB()
	repo := supplier.NewRepository(db)
	service := supplier.NewService(repo)

	fmt.Println("\033[33m🔌 Conectando a NATS...\033[0m")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a NATS: %v", err)
	}
	fmt.Println("\033[32m✅ Conexión a NATS exitosa.\033[0m")

	natsH := supplier.NewHandler(service)
	natsH.Subscribe(nc)

	select {} // Prevent main from exiting
}
