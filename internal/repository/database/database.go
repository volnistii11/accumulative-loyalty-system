package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) RegisterUser(user *model.User) error {
	if result := s.db.Create(user); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetUser(user *model.User) *model.User {
	s.db.Where(&user).Find(&user)
	return user
}

func (s *Storage) AddOrder(accumulation *model.Accumulation) error {
	if result := s.db.Select("user_id", "order_number", "uploaded_at", "processing_status").Create(accumulation); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAllOrders(userID int) *model.Accumulations {
	var orders model.Accumulations
	s.db.Select("order_number", "uploaded_at", "processing_status", "amount").Where("user_id = ?", userID).Find(&orders)
	return &orders
}

func (s *Storage) GetUserBalance(userID int) *model.Balance {
	var balance model.Balance
	s.db.Table("accumulations").Select("SUM(amount) as current").Where("user_id = ? AND amount > 0", userID).Find(&balance)
	s.db.Table("accumulations").Select("SUM(amount) as withdrawn").Where("user_id = ? AND amount < 0", userID).Find(&balance)
	return &balance
}

func (s *Storage) Withdraw(accumulation *model.Accumulation) error {
	if result := s.db.Select("user_id", "order_number", "uploaded_at", "amount").Create(accumulation); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAllUserWithdrawals(userID int) (*model.Withdrawals, error) {
	var (
		withdrawals *model.Withdrawals
		err         error
	)
	// TODO: get all users withdrawals
	return withdrawals, err
}

func (s *Storage) OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool {
	result := s.db.Where(&accumulation).Find(&accumulation)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func (s *Storage) OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool {
	result := s.db.
		Where("user_id != ? AND order_number = ?", accumulation.UserID, accumulation.OrderNumber).
		Find(accumulation)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func RunMigrations(dsn string) error {
	const migrationsPath = "./migrations"

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		return errors.Wrap(err, "start migrations")
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return errors.Wrap(err, "run migrations")
		}
	}
	return nil
}
