package repository

import (
	"database/sql"
	"fmt"
	"short_url/short/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertIfNotExists_PanicDBQuery(t *testing.T) {
	repo, db, mock := initRepoDbMock(t)
	defer db.Close()
	
	queryErr := fmt.Errorf("query error")
	mock.ExpectQuery(query).
	WithArgs("1234567").
	WillReturnError(queryErr)

	configToConnect(db)
	assert.PanicsWithError(t, queryErr.Error(), func(){repo.InsertIfNotExists("1234567", model.Url{})})

}

func TestInsertIfNotExists_PanicRowScan(t *testing.T) {
	repo, db, mock := initRepoDbMock(t)
	defer db.Close()

	rowErr := fmt.Errorf("read row error")
	rows := sqlmock.NewRows([]string{"LONG_URL", "TINY_URL"})
	rows.AddRow(0, "long.url.com").
	RowError(0, rowErr)

	mock.ExpectQuery(query).
	WithArgs("1234567").
	WillReturnRows(rows)

	configToConnect(db)
	assert.PanicsWithError(t, rowErr.Error(), func(){repo.InsertIfNotExists("1234567", model.Url{})})
}

func TestInsertIfNotExists_PanicPrepareContext(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abcd"
	long := "long.com/long"
	user_id := "user_id"
	err := fmt.Errorf("insert error")

	rows := sqlmock.NewRows([]string{"LONG_URL", "TINY_URL"})
	mock.ExpectQuery(query).
	WithArgs(tiny).
	WillReturnRows(rows)

	mock.ExpectPrepare(insert).
	ExpectExec().
	WithArgs(tiny,long,user_id).
	WillDelayFor(time.Second).
	WillReturnError(err)


	configToConnect(mockDB)
	assert.PanicsWithError(t, err.Error(), func(){repo.InsertIfNotExists(tiny, model.Url{Long:long, User:user_id})})

}

func TestInsertIfNotExists_PanicExecContext (t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abcd"
	long := "long.com/long"
	user_id := "user_id"
	err := fmt.Errorf("select error")

	rows := sqlmock.NewRows([]string{"LONG_URL", "TINY_URL"})
	mock.ExpectQuery(query).
	WithArgs(tiny).
	WillReturnRows(rows)

	mock.ExpectPrepare(insert)

	mock.ExpectExec(insert).
	WithArgs(tiny,long,user_id).
	WillDelayFor(time.Second).
	WillReturnError(err)

	configToConnect(mockDB)
	assert.PanicsWithError(t, err.Error(), func(){repo.InsertIfNotExists(tiny, model.Url{Long:long, User:user_id})})
}

func TestInsertIfNotExists_PanicRowsAffected(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abcd"
	long := "long.com/long"
	user_id := "user_id"
	err := fmt.Errorf("row affected error")

	rows := sqlmock.NewRows([]string{"LONG_URL", "TINY_URL"})
	mock.ExpectQuery(query).
	WithArgs(tiny).
	WillReturnRows(rows)

	mock.ExpectPrepare(insert)

	rowsAffectedErr := sqlmock.NewErrorResult(err)

	mock.ExpectExec(insert).
	WithArgs(tiny,long,user_id).
	WillDelayFor(time.Second).
	WillReturnResult(rowsAffectedErr)

	configToConnect(mockDB)
	assert.PanicsWithError(t, err.Error(), func(){repo.InsertIfNotExists(tiny, model.Url{Long:long, User:user_id})})
}

func TestInsertIfNotExists_RowsAffectedZero(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abcd"
	long := "long.com/long"
	user_id := "user_id"

	rows := sqlmock.NewRows([]string{"LONG_URL", "TINY_URL"})
	mock.ExpectQuery(query).
	WithArgs(tiny).
	WillReturnRows(rows)

	mock.ExpectPrepare(insert)

	mock.ExpectExec(insert).
	WithArgs(tiny,long,user_id).
	WillDelayFor(time.Second).
	WillReturnResult(sqlmock.NewResult(0,0))

	configToConnect(mockDB)
	assert.Equal(t, false, repo.InsertIfNotExists(tiny, model.Url{Long:long, User:user_id}))
}

func TestInsertIfNotExists_RowsAffectedOne(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abcd"
	long := "long.com/long"
	user_id := "user_id"

	rows := sqlmock.NewRows([]string{"LONG_URL", "TINY_URL"})
	mock.ExpectQuery(query).
	WithArgs(tiny).
	WillReturnRows(rows)

	mock.ExpectPrepare(insert)

	mock.ExpectExec(insert).
	WithArgs(tiny,long,user_id).
	WillDelayFor(time.Second).
	WillReturnResult(sqlmock.NewResult(0,1))

	configToConnect(mockDB)
	assert.Equal(t, true, repo.InsertIfNotExists(tiny, model.Url{Long:long, User:user_id}))
}

func TestRead_PanicQuery(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abc"
	err := fmt.Errorf("select err")

	mock.ExpectQuery(query_all).
	WithArgs(tiny).
	WillReturnError(err)

	configToConnect(mockDB)
	assert.PanicsWithError(t, err.Error(), func(){ repo.Read(tiny) })
}

func TestRead_EmptyResult(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abc"
	rows := sqlmock.NewRows([]string{})
	
	mock.ExpectQuery(query_all).
	WithArgs(tiny).
	WillReturnRows(rows)

	configToConnect(mockDB)
	assert.Empty(t, repo.Read(tiny))
}

func TestRead_Result(t *testing.T) {
	repo, mockDB, mock := initRepoDbMock(t)
	defer mockDB.Close()

	tiny := "abc"
	rows := sqlmock.
	NewRows([]string{"LONG_URL", "TINY_URL"}).
	AddRow("facebook.com","user_id")
	
	mock.ExpectQuery(query_all).
	WithArgs(tiny).
	WillReturnRows(rows)

	configToConnect(mockDB)
	assert.Equal(t, model.Url{Long:"facebook.com", User:"user_id"}, repo.Read(tiny))
}

func configToConnect(db *sql.DB) {
	toConnect = func(datasource DataSource) *sql.DB {
		return db
	}
}

func initRepoDbMock(t *testing.T) (Repository, *sql.DB, sqlmock.Sqlmock) {
	repo := New(DataSource{})
	// db, mock, err := sqlmock.New()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	processError(err, t)
	return repo, db, mock
}

func processError(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("error %s databases",err)
	}
}
