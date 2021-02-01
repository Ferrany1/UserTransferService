package db_service

import (
	"UserTransferService/src/system/config"
	"UserTransferService/src/system/l2f"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// config string for a connection
const dbConnString = "host=%s user=%s password=%s dbname=%s port=%d"

type Database struct {
	DB *gorm.DB
}

type User struct {
	Username string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
}

type Balance struct {
	Username string `gorm:"column:username" json:"username"`
	Balance  int    `gorm:"column:balance" json:"balance"`
}

// DB instance for all modules
var DB Database

// DB connector
func Connect() (err error) {
	newLogger := logger.New(
		l2f.Log, // io writer
		logger.Config{
			SlowThreshold: 400 * time.Second, // Slow SQL threshold
			LogLevel:      logger.Silent,     // Log level
			Colorful:      false,             // Disable color
		},
	)
	DB.DB, err = gorm.Open(postgres.Open(formatSourceName()), &gorm.Config{Logger: newLogger})
	return err
}

// Formats dbConnString from config
func formatSourceName() string {
	return fmt.Sprintf(dbConnString, config.CF.DB.Domain, config.CF.DB.Username, config.CF.DB.Password, config.CF.DB.DBName, config.CF.DB.Port)
}

// Gets user from users table
func (d Database) GetUser(u User) (user User, err error) {
	if result := d.DB.Table(config.CF.DB.Tables.Users).Where(&u).Find(&user); result.Error != nil {
		return user, errors.New(fmt.Sprintf("for GetUser failed to get Query resp %s", result.Error.Error()))
	}
	return
}

// Gets users balance from balance table
func (d Database) GetBalance(user User) (bal Balance, err error) {
	if result := d.DB.Table(config.CF.DB.Tables.Balances).Where(Balance{Username: user.Username}).Find(&bal); result.Error != nil {
		return bal, errors.New(fmt.Sprintf("for GetBalance failed to get Query resp %s", result.Error.Error()))
	}
	return
}

// Creates user's instance and balance
func (d Database) CreateUser(user User) (err error) {
	if result := d.DB.Table(config.CF.DB.Tables.Users).Create(&user); result.Error != nil {
		return errors.New(fmt.Sprintf("for CreateUser (user) failed to get Query resp %s", result.Error.Error()))
	}
	if result := d.DB.Table(config.CF.DB.Tables.Balances).Create(Balance{Username: user.Username, Balance: 0}); result.Error != nil {
		return errors.New(fmt.Sprintf("for CreateUser (balance) failed to get Query resp %s", result.Error.Error()))
	}
	return
}

// Gets all jobs from jobs table where working isn't a 1
func (d Database) UpdateUser(user User) (err error) {
	if result := d.DB.Table(config.CF.DB.Tables.Users).Where(User{Username: user.Username}).Updates(&user); result.Error != nil {
		return errors.New(fmt.Sprintf("for UpdateUser failed to get Query resp %s", result.Error.Error()))
	}
	return
}

// Makes updated transaction for sender and receiver balances
func (d Database) UpdateBalance(senderBal, receiverBal Balance) (err error) {
	d.DB.Transaction(
		func(tx *gorm.DB) (err error) {
			if result := tx.Table(config.CF.DB.Tables.Balances).Where(Balance{Username: senderBal.Username}).Updates(&senderBal); result.Error != nil {
				return errors.New(fmt.Sprintf("for UpdateBalance failed to get Query resp %s", result.Error.Error()))
			}
			tx.Transaction(func(tx2 *gorm.DB) (err error) {
				if result := tx2.Table(config.CF.DB.Tables.Balances).Where(Balance{Username: receiverBal.Username}).Updates(&receiverBal); result.Error != nil {
					return errors.New(fmt.Sprintf("for UpdateBalance failed to get Query resp %s", result.Error.Error()))
				}
				return
			})
			return
		},
		)
	return
}

// Deletes user instance and balance
func (d Database) DeleteUser(user User) (err error) {
	if result := d.DB.Table(config.CF.DB.Tables.Users).Where(User{Username: user.Username}).Delete(&user); result.Error != nil {
		return errors.New(fmt.Sprintf("for DeleteUser (user) failed to get Query resp %s", result.Error.Error()))
	}
	if result := d.DB.Table(config.CF.DB.Tables.Balances).Where(Balance{Username: user.Username}).Delete(Balance{Username: user.Username}); result.Error != nil {
		return errors.New(fmt.Sprintf("for DeleteUser (balance) failed to get Query resp %s", result.Error.Error()))
	}
	return
}
