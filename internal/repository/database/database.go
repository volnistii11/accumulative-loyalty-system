package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/cerrors"
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
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if s.OrderExistsAndBelongsToTheUser(accumulation) {
			return cerrors.ErrDBOrderExistsAndBelongsToTheUser
		}

		if s.OrderExistsAndDoesNotBelongToTheUser(accumulation) {
			return cerrors.ErrDBOrderExistsAndDoesNotBelongToTheUser
		}

		if result := tx.Select("user_id", "order_number", "uploaded_at", "processing_status").Create(accumulation); result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetAllOrders(userID int) ([]model.Accumulation, error) {
	var orders []model.Accumulation
	result := s.db.Select("order_number", "uploaded_at", "processing_status", "amount").Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (s *Storage) GetUserBalance(userID int) *model.Balance {
	var balance model.Balance
	s.db.Table("accumulations").Select("SUM(amount) as current").Where("user_id = ? AND amount > 0", userID).Find(&balance)
	s.db.Table("accumulations").Select("SUM(amount) as withdrawn").Where("user_id = ? AND amount < 0", userID).Find(&balance)
	return &balance
}

func (s *Storage) Withdraw(accumulation *model.Accumulation) error {
	if result := s.db.Select("user_id", "order_number", "processed_at", "amount").Create(accumulation); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAllUserWithdrawals(userID int) *model.Withdrawals {
	var withdrawals model.Withdrawals

	s.db.Table("accumulations").
		Select("order_number", "ABS(amount) as amount", "processed_at").
		Where("user_id = ? AND amount < 0", userID).
		Order("processed_at").Find(&withdrawals)

	return &withdrawals
}

func (s *Storage) OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool {
	result := s.db.Where(&accumulation).Find(&accumulation)
	fmt.Println(result)
	return result.RowsAffected > 0
}

func (s *Storage) OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool {
	result := s.db.
		Where("user_id != ? AND order_number = ?", accumulation.UserID, accumulation.OrderNumber).
		Find(accumulation)
	return result.RowsAffected > 0
}

func (s *Storage) GetNewOrders() []string {
	var orders []string
	s.db.Table("accumulations").
		Select("order_number").
		Where("processing_status = 'NEW' OR processing_status = 'PROCESSING'").
		Find(&orders)
	return orders
}

func (s *Storage) UpdateAccrualInfoForOrderNumber(newInfo *model.AccrualSystemAnswer) error {
	var accumulation model.Accumulation
	result := s.db.Model(&accumulation).
		Where("order_number = ?", newInfo.OrderNumber).
		Update("processing_status", newInfo.AccrualStatus).
		Update("accrual_status", newInfo.AccrualStatus).
		Update("amount", newInfo.Amount)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
