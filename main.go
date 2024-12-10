package main

import (
	"golang-app/handlers"
	"golang-app/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	dsn := "host=db user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	models.AutoMigrate(db)
	SeedUsers(db)
	SeedKTPRequests(db)
	SeedKartuKeluarga(db)

	r := gin.Default()

	r.POST("/login", func(c *gin.Context) { handlers.Login(c, db) })
	r.POST("/register", func(c *gin.Context) { handlers.CreateUser(c, db) })

	api := r.Group("/api")
	{
		api.Use(handlers.AuthMiddleware(""))                                                    // Ensure authentication middleware is applied
		api.GET("/list-permohonan", func(c *gin.Context) { handlers.GetAllKTPRequests(c, db) }) // Add this line for the new endpoint
		api.GET("/kk/:noKK", func(c *gin.Context) { handlers.GetKKByKK(c, db) })                // Add this line for the new endpoint
		api.GET("/permohonan", func(c *gin.Context) { handlers.GetKTPRequests(c, db) })         // Add this line for the new endpoint
		api.POST("/permohonan", func(c *gin.Context) { handlers.CreateKTPRequest(c, db) })      // Add this line for the new endpoint

		api.GET("/tracker", func(c *gin.Context) { handlers.GetUserRequestTracker(c, db) }) // Add this line for the new endpoint
		api.POST("/verify", func(c *gin.Context) { handlers.Verify(c, db) })                // Add this line for the new endpoint

		api.POST("/ktp", func(c *gin.Context) { handlers.CreateKTP(c, db) }) // Add this line for the new endpoint
		api.GET("/ktp", func(c *gin.Context) { handlers.GetKTP(c, db) })     // Add this line for the new endpoint

	}

	r.Run(":9090")
}

func SeedUsers(db *gorm.DB) {
	users := []models.User{
		{Name: "pak rt", Email: "rt@mail.com", Password: "123456", Role: "rt"},
		{Name: "pak w2", Email: "rw@mail.com", Password: "123456", Role: "rw"},
		{Name: "pak lurah", Email: "kelurahan@mail.com", Password: "123456", Role: "kelurahan"},
		{Name: "pak camat", Email: "kecamatan@mail.com", Password: "123456", Role: "kecamatan"},
		{Name: "atang", Email: "user@mail.com", Password: "123456", Role: "user"},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			panic("Failed to seed user: " + err.Error())
		}
	}
}

func SeedKTPRequests(db *gorm.DB) {
	ktpRequests := []models.KTPRequest{
		{NamaWarga: "John Doe", Email: "johndoe@mail.com", AlamatWarga: "Jl. Kebon Jeruk", NoKK: "1234567890", TglPengajuan: "2024-12-01", Status: models.StatusSubmitted},
		{NamaWarga: "Jane Smith", Email: "janesmith@mail.com", AlamatWarga: "Jl. Mangga Dua", NoKK: "2345678901", TglPengajuan: "2024-12-02", Status: models.StatusVerifiedRT},
		{NamaWarga: "Michael Johnson", Email: "mjohnson@mail.com", AlamatWarga: "Jl. Karet", NoKK: "3456789012", TglPengajuan: "2024-12-03", Status: models.StatusVerifiedRW},
	}

	for _, request := range ktpRequests {
		if err := db.FirstOrCreate(&request, models.KTPRequest{NamaWarga: request.NamaWarga}).Error; err != nil {
			panic("Failed to seed KTP request: " + err.Error())
		}
	}
}

func SeedKartuKeluarga(db *gorm.DB) {
	// Create sample data for Kartu Keluarga (KK)
	kartuKeluargas := []models.KartuKeluarga{
		{NoKK: "1234567890", NamaKK: "Budi Santoso", Alamat: "Jl. Kebon Jeruk", AnggotaKK: "Budi, Siti, Dwi"},
		{NoKK: "2345678901", NamaKK: "Ari Wibowo", Alamat: "Jl. Mangga Dua", AnggotaKK: "Ari, Rina, Wira"},
		{NoKK: "3456789012", NamaKK: "Mikael Pratama", Alamat: "Jl. Karet", AnggotaKK: "Mikael, Nita, Andi"},
		{NoKK: "4567890123", NamaKK: "Agus Santoso", Alamat: "Jl. Raya Bogor", AnggotaKK: "Agus, Tini, Rizki"},
		{NoKK: "5678901234", NamaKK: "Joko Widodo", Alamat: "Jl. Sudirman", AnggotaKK: "Joko, Siti, Ika"},
		{NoKK: "6789012345", NamaKK: "Yudi Handoko", Alamat: "Jl. Raya Serpong", AnggotaKK: "Yudi, Sari, Dita"},
		{NoKK: "7890123456", NamaKK: "Susi Susanti", Alamat: "Jl. Pahlawan", AnggotaKK: "Susi, Agus, Bayu"},
		{NoKK: "8901234567", NamaKK: "Anton Setiawan", Alamat: "Jl. Mataram", AnggotaKK: "Anton, Lina, Eka"},
		{NoKK: "9012345678", NamaKK: "Rudi Hartono", Alamat: "Jl. Cendrawasih", AnggotaKK: "Rudi, Lila, Dede"},
		{NoKK: "0123456789", NamaKK: "Tono Prakoso", Alamat: "Jl. Merdeka", AnggotaKK: "Tono, Maria, Kevin"},
	}

	// Loop through and insert the sample data into the database
	for _, kk := range kartuKeluargas {
		if err := db.FirstOrCreate(&kk, models.KartuKeluarga{NoKK: kk.NoKK}).Error; err != nil {
			panic("Failed to seed Kartu Keluarga" + err.Error())
		}
	}

	log.Println("Successfully seeded 10 Kartu Keluarga records.")
}
