package database

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ItemPrices = make(map[string]int, 10)

type Item struct {
	Name  string `json:"name" gorm:"primaryKey"`
	Price int    `json:"price"`
}

type User struct {
	Username  string `json:"username" gorm:"uniqueIndex"`
	Inventory []InventoryItem
	ID        uint `json:"id" gorm:"primaryKey"`
	Balance   int  `json:"balance" gorm:"default:1000"`
}

type Transaction struct {
	CreatedAt time.Time
	ID        uint `gorm:"primaryKey"`
	FromUser  uint
	ToUser    uint
	Amount    int
}

type InventoryItem struct {
	ItemType string
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	Quantity int
}

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if ok := DB.Migrator().HasTable("items"); !ok {
		createTableSQL := `CREATE TABLE IF NOT EXISTS items(name varchar PRIMARY KEY, price INT NOT NULL)`

		err := DB.Exec(createTableSQL).Error
		if err != nil {
			log.Fatal("Ошибка создания таблицы:", err)
		}

		file, err := os.Open("items.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		log.Println("ПУПУПУ")
		for scanner.Scan() {
			fields := strings.Split(scanner.Text(), ",")
			DB.Exec("INSERT INTO items (name, price) VALUES ($1, $2)", fields[0], fields[1])
			log.Println(fields[0], fields[1])
		}
	}

	rows, err := DB.Raw(`SELECT name, price FROM items;`).Rows()
	if err != nil {
		log.Fatal("БД items не найдена")
	}
	defer rows.Close()

	for rows.Next() {
		var price int
		var name string
		if err := rows.Scan(&name, &price); err != nil {
			log.Fatal("ошибка в БД items")
		}

		ItemPrices[name] = price
	}
	if err := DB.AutoMigrate(&User{}, &Transaction{}, &InventoryItem{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
}
