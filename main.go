package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n==========================================\n", sql)
}

var db *gorm.DB

func main() {
	dsn := "root:P@ssw0rd@tcp(localhost:3306)/gorm_test_inf?parseTime=true"
	dial := mysql.Open(dsn)

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false, //ไม่ทำจริงใน db ถ้า true
	})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(Gender{}, Test{}, Customer{})
	// CreateGender("xxxx")
	// GetGenders()
	// GetGender(1)
	// GetGenderByName("Male")
	// UpdateGender2(4, "")
	// DeleteGender(4)

	// db.Migrator().CreateIndex(Customer{})
	// db.Migrator().CreateTable(Customer{})

	// CreateCustomer("Note", 2)
	GetCustomers()
}

func GetCustomers() {
	customers := []Customer{}
	// tx := db.Preload("Gender").Find(&customers)
	tx := db.Preload(clause.Associations).Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	// fmt.Println(customers)
	for _, customer := range customers {
		fmt.Printf("%v|%v|%v\n", customer.ID, customer.Name, customer.Gender.Name)
	}

}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{Name: name, GenderId: genderID}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

type Customer struct {
	ID       uint
	Name     string
	Gender   Gender
	GenderId uint
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("deleted")
	GetGender(id)
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
}

func UpdateGender2(id uint, name string) { //ค่าต้องไม่เป็น zero value ไม่งั้นจะไม่ทำงาน
	gender := Gender{Name: name}
	// tx := db.Model(&Gender{}).Where("id = ?", id).Updates(gender)
	tx := db.Model(&Gender{}).Where("id = @myid", sql.Named("myid", id)).Updates(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	GetGender(id)
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGender(id uint) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByName(name string) {
	gender := Gender{}
	tx := db.First(&gender, "name = ?", name)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func CreateGender(name string) {
	gender := Gender{Name: name}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	fmt.Println(gender)
}

type Gender struct {
	ID   uint
	Name string `gorm:"unique;size(10)"`
}

type Test struct {
	gorm.Model
	Code uint   `gorm:"primaryKey;comment:This is Code"`
	Name string `gorm:"column:myname;size:20;unique;default:Hello;not null"`
}

func (t Test) TableName() string {
	return "MyTest"
}
