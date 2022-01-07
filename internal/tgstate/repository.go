package tgstate

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type stateDB struct { // nolint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	TelegramID int   `db:"telegram_id"`
	State      State `db:"state"`
}

type Repository interface {
	Get(telegramID int) (*State, error)
	GetAll() ([]*State, error)
	Save(state *State) error
}

type Repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

// Get returns a state by telegram ID.
// If there is no state in the database, a new one will be created.
func (r *Repo) Get(telegramID int) (*State, error) {
	row := new(stateDB)
	err := r.db.
		QueryRow("SELECT state FROM states WHERE telegram_id = $1", telegramID).
		Scan(&row.State)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		return &row.State, nil
	}

	// There is no state with such telegram id. Create a new record.
	row = &stateDB{
		CreatedAt:  time.Now(),
		TelegramID: telegramID,
		State: State{
			TelegramID: telegramID,
		},
	}
	row.UpdatedAt = row.CreatedAt

	query := `INSERT INTO states(created_at, updated_at, telegram_id, state) 
		VALUES(:created_at, :updated_at, :telegram_id, :state)`
	_, err = r.db.NamedExec(query, row)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new state")
	}

	return &row.State, nil
}

// GetAll fetches all states and returns them.
func (r *Repo) GetAll() ([]*State, error) {
	var records []stateDB
	err := r.db.Select(&records, `SELECT * FROM states`)
	if err != nil {
		return nil, err
	}

	result := make([]*State, len(records))
	for i := range records {
		result[i] = &records[i].State
	}

	return result, nil
}

// Save saves the given state to the database.
func (r *Repo) Save(state *State) error {
	row := &stateDB{
		UpdatedAt:  time.Now(),
		TelegramID: state.TelegramID,
		State:      *state,
	}
	_, err := r.db.NamedExec(`UPDATE states 
		SET updated_at=:updated_at, state=:state 
		WHERE telegram_id = :telegram_id`,
		row)

	return err
}
