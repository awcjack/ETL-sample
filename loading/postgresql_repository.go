package loading

import (
	"context"
	"time"

	"github.com/awcjack/ETL-sample/config"
	"github.com/awcjack/ETL-sample/transformation"
	"github.com/awcjack/ETL-sample/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQLRepository struct {
	db     *sqlx.DB
	logger utils.Logger
}

func NewPostgreSQLRepository(db *sqlx.DB, logger utils.Logger) *PostgreSQLRepository {
	if db == nil {
		logger.Panicf("missing db")
	}

	return &PostgreSQLRepository{
		db:     db,
		logger: logger,
	}
}

type postgresqlUser struct {
	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	DateOfBirth   time.Time `db:"date_of_birth"`
	City          string    `db:"city"`
	StreetName    string    `db:"street_name"`
	StreetAddress string    `db:"street_address"`
	ZipCode       string    `db:"zip_code"`
	State         string    `db:"state"`
	Country       string    `db:"country"`
	Latitude      float64   `db:"latitude"`
	Longitude     float64   `db:"longitude"`
}

// insert
func (p *PostgreSQLRepository) AddUser(ctx context.Context, user transformation.TransformedData) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	defer func() {
		err = p.finishTransaction(err, tx)
	}()

	dbUser := postgresqlUser{
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		DateOfBirth:   user.DateOfBirth,
		City:          user.Address.City,
		StreetName:    user.Address.StreetName,
		StreetAddress: user.Address.StreetAddress,
		ZipCode:       user.Address.ZipCode,
		State:         user.Address.State,
		Country:       user.Address.Country,
		Latitude:      user.Address.Latitude,
		Longitude:     user.Address.Longitude,
	}

	// Insert user to users table in PostgreSQL
	_, err = tx.NamedExec(`
		INSERT INTO
			users (first_name, last_name, date_of_birth, city, street_name, street_address, zip_code, state, country, latitude, longitude)
		VALUES
			(:first_name, :last_name, :date_of_birth, :city, :street_name, :street_address, :zip_code, :state, :country, :latitude, :longitude)
	`, dbUser)
	if err != nil {
		return err
	}

	return nil
}

// bulk insert
func (p *PostgreSQLRepository) AddUsers(ctx context.Context, users []transformation.TransformedData) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	defer func() {
		err = p.finishTransaction(err, tx)
	}()

	dbUsers := make([]postgresqlUser, 0, len(users))
	for _, user := range users {
		dbUsers = append(dbUsers, postgresqlUser{
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			DateOfBirth:   user.DateOfBirth,
			City:          user.Address.City,
			StreetName:    user.Address.StreetName,
			StreetAddress: user.Address.StreetAddress,
			ZipCode:       user.Address.ZipCode,
			State:         user.Address.State,
			Country:       user.Address.Country,
			Latitude:      user.Address.Latitude,
			Longitude:     user.Address.Longitude,
		})
	}

	// Insert user to users table in PostgreSQL
	result, err := tx.NamedExec(`
		INSERT INTO
			users (first_name, last_name, date_of_birth, city, street_name, street_address, zip_code, state, country, latitude, longitude)
		VALUES
			(:first_name, :last_name, :date_of_birth, :city, :street_name, :street_address, :zip_code, :state, :country, :latitude, :longitude)
	`, dbUsers)
	p.logger.Infof("result %v", result)
	p.logger.Infof("err %v", err)
	if err != nil {
		return err
	}

	return nil
}

// PostgreSQL finish transaction operation (rollback if failure and commit if no error)
func (p *PostgreSQLRepository) finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrap(err, rbErr.Error())
		}

		return err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return err
}

// Start PostgreSQL connection
func NewPostgreSQLConnection(c config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", c.ConnectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
