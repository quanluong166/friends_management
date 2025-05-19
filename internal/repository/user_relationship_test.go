package repository_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quanluong166/friends_management/internal/constant"
	"github.com/quanluong166/friends_management/internal/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		_ = db.Close()
	}

	return gdb, mock, cleanup
}

func TestCreateFriendRelationship_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	email1 := "alice@example.com"
	email2 := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(email1, email2, constant.FRIEND_RELATIONSHIP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	err := repo.CreateFriendRelationship(email1, email2)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateFriendRelationship_FailCreateFirstRelation(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	email1 := "alice@example.com"
	email2 := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(email1, email2, constant.FRIEND_RELATIONSHIP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(email2, email1, constant.FRIEND_RELATIONSHIP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectRollback()

	tx := db.Begin()
	err := repo.CreateFriendRelationship(email1, email2)
	require.Error(t, err)

	err = tx.Commit().Error
	require.Error(t, err)
	require.Error(t, mock.ExpectationsWereMet())

}

func TestCreateFriendRelationship_FailCreateSecondRelation(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	email1 := "alice@example.com"
	email2 := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(email1, email2, constant.FRIEND_RELATIONSHIP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(email2, email1, constant.FRIEND_RELATIONSHIP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	tx := db.Begin()
	err := repo.CreateFriendRelationship(email1, email2)
	require.Error(t, err)

	err = tx.Commit().Error
	require.Error(t, err)
	require.Error(t, mock.ExpectationsWereMet())
}

func TestGetListSubscriberEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	targetEmail := "bob@example.com"

	rows := sqlmock.NewRows([]string{"requestor_email", "target_email", "type"}).
		AddRow("alice@example.com", targetEmail, constant.SUBSCRIBER_RELATIONSHIOP_TYPE).
		AddRow("john@example.com", targetEmail, constant.SUBSCRIBER_RELATIONSHIOP_TYPE)

	mock.ExpectQuery(`SELECT \* FROM "user_relationships"`).
		WithArgs(targetEmail, constant.SUBSCRIBER_RELATIONSHIOP_TYPE).
		WillReturnRows(rows)

	result, err := repo.GetListSubscriberEmail(targetEmail)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"alice@example.com", "john@example.com"}, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetListSubscriberEmail_FailDatabase(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	targetEmail := "bob@example.com"

	mock.ExpectQuery(`SELECT \* FROM "user_relationships"`).
		WithArgs(targetEmail, constant.SUBSCRIBER_RELATIONSHIOP_TYPE).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.GetListSubscriberEmail(targetEmail)
	require.Error(t, err)
	require.Nil(t, nil, result)
}

func TestGetListFriendshipEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestorEmail := "test1@example.com"

	rows := sqlmock.NewRows([]string{"requestor_email", "target_email", "type"}).
		AddRow(requestorEmail, "test2@example.com", constant.FRIEND_RELATIONSHIP_TYPE).
		AddRow(requestorEmail, "test3@example.com", constant.FRIEND_RELATIONSHIP_TYPE)

	mock.ExpectQuery(`SELECT \* FROM "user_relationships"`).
		WithArgs(requestorEmail, constant.FRIEND_RELATIONSHIP_TYPE).
		WillReturnRows(rows)

	result, err := repo.GetListFriendshipEmail(requestorEmail)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"test2@example.com", "test3@example.com"}, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckTwoUsersAreFriends(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	email1 := "email1@example.com"
	email2 := "email2@example.com"
	rows := sqlmock.NewRows([]string{"requestor_email", "target_email", "type"}).
		AddRow(email1, email2, constant.FRIEND_RELATIONSHIP_TYPE).
		AddRow(email2, email1, constant.FRIEND_RELATIONSHIP_TYPE)

	mock.ExpectQuery(`SELECT \* FROM "user_relationships"`).
		WithArgs(email1, email2, constant.FRIEND_RELATIONSHIP_TYPE).
		WillReturnRows(rows)

	isFriend, err := repo.CheckTwoUsersAreFriends(email1, email2)
	require.NoError(t, err)
	require.Equal(t, isFriend, true)

}

func TestCheckTwoUsersBlockedEachOther(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	email1 := "email1@example.com"
	email2 := "email2@example.com"
	rows := sqlmock.NewRows([]string{"requestor_email", "target_email", "type"}).
		AddRow(email1, email2, constant.BLOCK_RELATIONSHIP_TYPE).
		AddRow(email2, email1, constant.BLOCK_RELATIONSHIP_TYPE)

	mock.ExpectQuery(`SELECT \* FROM "user_relationships"`).
		WithArgs(email1, email2, constant.BLOCK_RELATIONSHIP_TYPE, email2, email1, constant.BLOCK_RELATIONSHIP_TYPE).
		WillReturnRows(rows)

	isBlock, err := repo.CheckTwoUsersBlockedEachOther(email1, email2)
	require.NoError(t, err)
	require.Equal(t, isBlock, true)
}

func TestAddSubscriber(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestor := "alice@example.com"
	target := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(requestor, target, constant.SUBSCRIBER_RELATIONSHIOP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.AddSubscriber(requestor, target)
	require.NoError(t, err)
}

func TestCreateBlockRelationship(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestor := "alice@example.com"
	target := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_relationships"`)).
		WithArgs(requestor, target, constant.BLOCK_RELATIONSHIP_TYPE, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateBlockRelationship(requestor, target)
	require.NoError(t, err)
}

func TestCheckIfTheRequestorAlreadySubscribe(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestor := "alice@example.com"
	target := "bob@example.com"

	rows := sqlmock.NewRows([]string{"requestor_email", "target_email", "type"}).
		AddRow(requestor, target, constant.SUBSCRIBER_RELATIONSHIOP_TYPE)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_relationships"`)).
		WithArgs(requestor, target, constant.SUBSCRIBER_RELATIONSHIOP_TYPE, sqlmock.AnyArg()).
		WillReturnRows(rows)

	isSubscriber, err := repo.CheckIfTheRequestorAlreadySubscribe(requestor, target)
	require.NoError(t, err)
	require.Equal(t, isSubscriber, true)
}

func TestDeleteRelationship(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestor := "alice@example.com"
	target := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_relationships"`)).
		WithArgs(requestor, target).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteRelationship(requestor, target)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRelationship_FailDeleteFirstRelationship(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestor := "alice@example.com"
	target := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_relationships"`)).
		WithArgs(target, requestor).
		WillReturnError(sql.ErrTxDone)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_relationships"`)).
		WithArgs(requestor, target).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()

	err := repo.DeleteRelationship(requestor, target)
	require.Error(t, err)
	require.Error(t, mock.ExpectationsWereMet())
}

func TestDeleteRelationship_FailDeleteSecondRelationship(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRelationshipRepository(db)

	requestor := "alice@example.com"
	target := "bob@example.com"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_relationships"`)).
		WithArgs(target, requestor).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_relationships"`)).
		WithArgs(requestor, target).
		WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx := db.Begin()
	err := repo.DeleteRelationship(requestor, target)
	require.Error(t, err)

	err = tx.Commit().Error
	require.Error(t, err)
	require.Error(t, mock.ExpectationsWereMet())
}
