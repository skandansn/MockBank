package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectDataBase() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		fmt.Println("We are getting the env values")
	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	DB, err = gorm.Open(Dbdriver, DBURL)

	if err != nil {
		fmt.Println("Cannot connect to database ", Dbdriver)
		log.Fatal("connection error:", err)
	} else {
		fmt.Println("We are connected to the database ", Dbdriver)
	}

	DB.AutoMigrate(&Customer{}, &User{}, &Employee{}, &Access{}, &Appointment{}, &BookedAppointment{}, &BankAccount{}, &Card{}, &CardRequest{}, &Transaction{}, &BankAccountRequest{})
	DB.Model(&Customer{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Employee{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Access{}).AddForeignKey("employee_id", "employees(id)", "RESTRICT", "RESTRICT")

	DB.Model(&Appointment{}).AddForeignKey("employee_id", "employees(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Appointment{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT")

	DB.Model(&BookedAppointment{}).AddForeignKey("employee_id", "employees(id)", "RESTRICT", "RESTRICT")
	DB.Model(&BookedAppointment{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT")
	DB.Model(&BookedAppointment{}).AddForeignKey("appointment_id", "appointments(id)", "RESTRICT", "RESTRICT")

	DB.Model(&BankAccount{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT")

	DB.Model(&Card{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT")

	DB.Model(&CardRequest{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT")

	DB.Model(&Transaction{}).AddForeignKey("sender_account", "bank_accounts(account_number)", "RESTRICT", "RESTRICT")
	DB.Model(&Transaction{}).AddForeignKey("receiver_account", "bank_accounts(account_number)", "RESTRICT", "RESTRICT")

	DB.Model(&BankAccountRequest{}).AddForeignKey("booked_appointment_id", "booked_appointments(id)", "RESTRICT", "RESTRICT")

}
