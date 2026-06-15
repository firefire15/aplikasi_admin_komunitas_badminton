package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"aplikasi_admin_komunitas_badminton/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct{
	CommunityID uint `json:"community_id" binding:"required"`
	PlayerID uint `json:"player_id" binding:"required"`
	Email		string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Role string `json:"role"` 

}

type LoginInput struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterAdmin(c *gin.Context){
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Gagal mengonversi password"})
		return
	}

	newAdmin := db.Admin{
		CommunityID: input.CommunityID,
		PlayerID: input.PlayerID,
		Email: input.Email,
		Password: string(hashedPassword),
		Role: input.Role,
	}

	if err := db.DB.Create(&newAdmin).Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar atau Community ID salah"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Admin berhasil didaftarkan!"})
}

func LoginAdmin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin db.Admin
	// Cari admin berdasarkan email
	if err := db.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Verifikasi password biner bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token, err := helper.GenerateJWT(admin.Email, admin.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   token,
	})
}

func CreateAdmin(c *gin.Context) {
	var input db.Admin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&input).Error; err != nil {
		// Menangani error jika email duplikat atau community_id tidak ada
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mendaftarkan admin. Cek kembali email atau Community ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin berhasil didaftarkan", "data": input})
}

func GetAdmins(c *gin.Context) {
	var listAdmins []db.Admin
	
	// Query SELECT * FROM admins beserta relasi data Komunitasnya (Preload)
	if err := db.DB.Preload("Community").Find(&listAdmins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": listAdmins})
}

func UpdateAdmin(c *gin.Context) {
	var input db.Admin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingAdmin db.Admin
	if err := db.DB.First(&existingAdmin, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	existingAdmin.Email = input.Email
	existingAdmin.Role = input.Role

	db.DB.Save(&existingAdmin)
	c.JSON(http.StatusOK, gin.H{"message": "Data admin berhasil diperbarui", "data": existingAdmin})
}