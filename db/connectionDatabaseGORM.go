package db

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Admin struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CommunityID uint           `gorm:"not null" json:"community_id"` 
	Community   Community      `gorm:"foreignKey:CommunityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"community,omitempty"`
	PlayerID    uint	       `gorm:"not null" json:"player_id"` 
	Player		Player         `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"community,omitempty"`
	Email       string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	Role        string         `gorm:"type:varchar(50);default:'Admin'" json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Community struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Location  string         `gorm:"type:text" json:"location"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` 
}

type Player struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone       string         `gorm:"type:varchar(20)" json:"phone"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CommunityPlayer struct{
	ID          uint           `gorm:"primaryKey" json:"id"`
	CommunityID uint           `gorm:"not null" json:"community_id"`
	Community   Community	   `gorm:"foreignKey:CommunityID" json:"community"`
	PlayerID	uint		   `gorm:"not null" json:"player_id"`
	Player      Player         `gorm:"foreignKey:PlayerID" json:"player"`
	Status      string		   `gorm:"type:varchar(20)";default:'active' json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Court struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CommunityID uint           `gorm:"not null" json:"community_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"` 
	PricePerDay float64        `gorm:"type:numeric" json:"price_per_day"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Schedule struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CourtID     uint           `gorm:"not null" json:"court_id"`
	Court       Court          `gorm:"foreignKey:CourtID" json:"court,omitempty"`
	Date        time.Time      `gorm:"type:date" json:"date"`
	StartTime   string         `gorm:"type:varchar(5)" json:"start_time"` 
	EndTime     string         `gorm:"type:varchar(5)" json:"end_time"`  
	Status      string         `gorm:"type:varchar(50);default:'Booked'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Shuttlecock struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CommunityID uint           `gorm:"not null" json:"community_id"`
	Brand       string         `gorm:"type:varchar(100);not null" json:"brand"` 
	Stock       int            `gorm:"default:0" json:"stock"`
	Returned    int			   `gorm:"default:0" json:"returned"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Receipt struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CommunityID uint           `gorm:"not null" json:"community_id"`
	Category    string         `gorm:"type:varchar(100)" json:"category"`
	Amount      float64        `gorm:"type:numeric;not null" json:"amount"`
	Notes       string         `gorm:"type:text" json:"notes"`
	Date        time.Time      `json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
}

type Match struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ScheduleID     uint           `gorm:"not null" json:"schedule_id"`
	FormatScore    string         `gorm:"type:varchar(50);default:'21 Rally'" json:"format_score"`
	CocksUsed      int            `gorm:"default:0" json:"cocks_used"`
	ShuttlecockID  uint           `json:"shuttlecock_id"`
	DurationMinute int            `json:"duration_minute"`
	ScoreSet1      string         `gorm:"type:varchar(10)" json:"score_set1"`
	ScoreSet2      string         `gorm:"type:varchar(10)" json:"score_set2"`
	ScoreSet3      string         `gorm:"type:varchar(10)" json:"score_set3"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CommunityPlayer   []CommunityPlayer  `gorm:"many2many:match_players;" json:"community_players"`
}

type Payment struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ScheduleID     uint           `gorm:"not null" json:"schedule_id"`
	CommunityID    uint           `gorm:"not null" json:"community_id"`
	CommunityPlayerID       uint           `gorm:"not null" json:"communityplayer_id"`
	CommunityPlayer         CommunityPlayer         `gorm:"foreignKey:CommunityPlayerID" json:"communityplayer"`
	CourtFee       float64        `gorm:"type:numeric" json:"court_fee"`       
	CockFee        float64        `gorm:"type:numeric" json:"cock_fee"`      
	TotalPaid      float64        `gorm:"type:numeric" json:"total_paid"`
	IsConfirmed    bool           `gorm:"default:false" json:"is_confirmed"`
	CreatedAt      time.Time      `json:"created_at"`
}

func ConnectDatabase() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat menemukan file .env")
	}

	host     := os.Getenv("PGHOST")
	port     := os.Getenv("PGPORT")
	user     := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbName   := os.Getenv("PGDATABASE")

	fmt.Println("data ", host, port, user, password, dbName)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	err = database.AutoMigrate(&Community{}, &Admin{}, &Player{}, &Court{},
		 &Schedule{}, &Shuttlecock{}, &Receipt{}, &Match{}, &Payment{}, &CommunityPlayer{})

	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}

	DB = database
	fmt.Println("Koneksi database berhasil dan migrasi selesai!")

}