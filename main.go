package main

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"aplikasi_admin_komunitas_badminton/db"
	"aplikasi_admin_komunitas_badminton/routes"
)

func main() {
	db.ConnectDatabase()

	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat menemukan file .env, sistem akan menggunakan env global")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server Bioskop berjalan di port :%s...\n", port)
	err = routes.StartAPIServer().Run(":" + port)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}