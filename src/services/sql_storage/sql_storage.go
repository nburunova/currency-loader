package sql_storage

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nburunova/currency-loader/src/api"
	"github.com/nburunova/currency-loader/src/currency"
	"github.com/nburunova/currency-loader/src/services/log"
	"github.com/pkg/errors"
)

const (
	dbName         = "currency"
	createTableSQL = `CREATE TABLE IF NOT EXISTS currency (id INTEGER PRIMARY KEY, code TEXT, value REAL, addtime INTEGER)`
	createViewSQL  = `CREATE VIEW IF NOT EXISTS latest_currency AS
	SELECT id, code, value
	FROM currency
	WHERE addtime = (SELECT MAX(addtime) FROM currency)
	ORDER BY id`
)

var (
	ErrEmptyPage   = errors.New("No data for page")
	ErrInvalidPage = errors.New("Invalid page")
)

type Storage struct {
	db       *sql.DB
	lifeTime time.Duration
	logger   log.Logger
}

func New(lifeTime time.Duration, logger log.Logger) (*Storage, error) {
	dbPath := fmt.Sprintf("./%v.db", dbName)
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create DB %v", dbPath)
	}
	stmt, err := database.Prepare(createTableSQL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot prepare currency table")
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create currency table")
	}
	stmt, err = database.Prepare(createViewSQL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot prepare latest currency view")
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create latest currency view")
	}
	return &Storage{
		db:       database,
		lifeTime: lifeTime,
		logger:   logger,
	}, nil
}

func (s *Storage) saveTs() int64 {
	return time.Now().UTC().Unix()
}

func (s *Storage) oldTs() int64 {
	return time.Now().Add(-s.lifeTime).UTC().Unix()
}

func (s *Storage) Save(ctx context.Context, valutes currency.Valutes) error {
	addTime := s.saveTs()
	sqlStr := "INSERT INTO currency(code, value, addtime) VALUES "
	vals := []interface{}{}

	for _, row := range valutes {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, strings.ToLower(string(row.Code())), row.Value(), addTime)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := s.db.Prepare(sqlStr)
	if err != nil {
		return errors.Wrap(err, "cannot prepare insert values to table")
	}
	//format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		return errors.Wrap(err, "cannot insert values to table")
	}
	return nil
}

func (s *Storage) TotalPages(ctx context.Context, pageSize int) (*int, error) {
	query := `SELECT count(*) FROM latest_currency`
	s.logger.Debug(query)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "cannot request total pages")
	}
	var counter float64
	for rows.Next() {
		err := rows.Scan(&counter)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan total pages value")
		}
	}
	counterInt := int(math.Ceil(counter / float64(pageSize)))
	return &counterInt, nil
}

func (s *Storage) ValutesOnPage(ctx context.Context, pageNum, pageSize int) (api.Valutes, error) {
	if pageNum <= 0 {
		return nil, errors.Wrapf(ErrInvalidPage, "%v", pageNum)
	}
	query := fmt.Sprintf(`
	SELECT code, value FROM latest_currency
	LIMIT %v, %v`, (pageNum-1)*pageSize, pageSize)
	s.logger.Debug(query)
	rows, err := s.db.Query(query)
	var code string
	var value float64
	if err != nil {
		return nil, errors.Wrap(err, "cannot request valutes")
	}
	result := make(api.Valutes, 0)
	for rows.Next() {
		err := rows.Scan(&code, &value)
		if err != nil {
			return nil, errors.Wrap(err, "cannot can code and value for valute")
		}
		result = append(result, api.NewValute(code, value))
	}
	if len(result) == 0 {
		return nil, ErrEmptyPage
	}
	return result, nil
}

func (s *Storage) ValuteValueByCode(ctx context.Context, code string) (*float64, error) {
	code = strings.ToLower(code)
	query := fmt.Sprintf(`
	SELECT value FROM latest_currency WHERE code = '%v'`, code)
	s.logger.Debug(query)
	rows, err := s.db.Query(query)
	var value float64
	if err != nil {
		return nil, errors.Wrapf(err, "cannot request valute value by code %v", code)
	}
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot scan value for code %v", code)
		}
	}
	return &value, nil
}

func (s *Storage) deletOld(ctx context.Context) error {
	result, err := s.db.Exec(fmt.Sprintf(`DELETE FROM currency WHERE addTime < %v`, s.oldTs()))
	if err != nil {
		return errors.Wrap(err, "cannot delete old valutes")
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get deleted valutes number")
	}
	s.logger.Debug(fmt.Sprintf("%v old valutes deleted", deleted))
	return nil
}

func (s *Storage) Start(ctx context.Context, checkPeriod time.Duration) {
	ticker := time.NewTicker(checkPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.deletOld(ctx)
				if err != nil {
					s.logger.Error(errors.Wrap(err, "cannot delete old valutes"))
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
