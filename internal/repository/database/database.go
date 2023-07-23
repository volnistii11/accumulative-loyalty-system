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

func (s *Storage) RegisterUser(login string, pass string) error {
	var (
		err error
	)
	// TODO: register user with some ORM
	return err
}

func (s *Storage) AuthenticateUser(login string, pass string) error {
	var (
		err error
	)
	// TODO: authenticate user with some ORM
	return err
}

func (s *Storage) PutOrder(orderNumber string) error {
	var (
		err error
	)
	// TODO: put in order number into database into accumulation table
	return err
}

func (s *Storage) GetAllOrders(userID string) (*model.Accumulations, error) {
	var (
		accumulations *model.Accumulations
		err           error
	)
	// TODO: get all orders by userID
	return accumulations, err
}

func (s *Storage) GetUserBalance(userID string) (float64, float64, error) {
	var (
		currentBalance   float64
		withdrawnBalance float64
		err              error
	)
	// TODO: get user balance by userID
	return currentBalance, withdrawnBalance, err
}

func (s *Storage) Withdraw(orderNumber int, amount int) error {
	var (
		err error
	)
	// TODO: withdraw from user balance
	return err
}

func (s *Storage) GetAllUserWithdrawals(userID int) (*model.Withdrawals, error) {
	var (
		withdrawals *model.Withdrawals
		err         error
	)
	// TODO: get all users withdrawals
	return withdrawals, err
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
