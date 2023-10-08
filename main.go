package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// object struct of transaction
type transaction struct {
	Id     int
	UserId int
	Date   time.Time
	Status string
	Amount int
}

// this function is return db variable for database integration
func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/db_financial?parseTime=true")
	if err != nil {
		log.Fatal("Connection failed")
	}

	return db
}

func main() {
listMenu:
	for {
		var selectedMenu int

		fmt.Println("Aplikasi pencatatan keuangan")
		fmt.Println("1. Lihat catatan")
		fmt.Println("2. Tambah catatan")
		fmt.Println("3. Update catatan")
		fmt.Println("4. Hapus catatan")
		fmt.Println("0. Exit")
		fmt.Print("Pilih (0-4) : ")
		fmt.Scan(&selectedMenu)

		switch selectedMenu {
		case 0:
			break listMenu
		case 1:
			GetAllTransactions()
		case 2:
			InsertNewTransaction()
		case 3:
			UpdateTransaction()
		case 4:
			DeleteTransaction()
		default:
			fmt.Println("Pilihan tidak valid")
			continue
		}
	}
}
