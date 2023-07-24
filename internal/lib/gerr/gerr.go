package gerr

import (
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

func GetCodeFromGORMMessage(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func IsDuplicateKey(err error) bool {
	if GetCodeFromGORMMessage(err) == pgerrcode.UniqueViolation {
		return true
	}
	return false
}
