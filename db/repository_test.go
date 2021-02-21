package db

//
//import (
//	"database/sql"
//	"database/sql/driver"
//	"reflect"
//	"regexp"
//	"testing"
//	"time"
//
//	"github.com/DATA-DOG/go-sqlmock"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//)
//
//var mock sqlmock.Sqlmock
//
//const sqlSelectAll = `SELECT * FROM "sets"`
//const sqlLatest = `SELECT * FROM "sets" WHERE (type in ($1,$2)) AND (release_date <= now()) ORDER BY release_date DESC,"sets"."code" ASC LIMIT 1`
//const sqlFindSet = `SELECT * FROM "sets" WHERE ( code = $1) ORDER BY "sets"."code" ASC LIMIT 1`
//const sqlFindAllCardsInSet = `SELECT * FROM "cards"  WHERE ("set_id" IN ($1)) ORDER BY "cards"."uuid" ASC`
//const sqlFindAllSheets = `SELECT * FROM "sheets"  WHERE ("set_id" IN ($1)) ORDER BY "sheets"."id" ASC`
//
//var setColumns = []string{"code", "name", "type", "release_date", "base_set_size", "pack_configurations"}
//var basicSetValues = []driver.Value{basicSet.Code, basicSet.Name, basicSet.Type, basicSet.ReleaseDate, basicSet.BaseSetSize, basicSet.PackConfigurations}
//
//var basicSet = Set{
//	Code:               "CODE",
//	Name:               "NAME",
//	Type:               "TYPE",
//	ReleaseDate:        time.Now(),
//	BaseSetSize:        0,
//	SheetCards:              nil,
//	Sheets:             nil,
//	PackConfigurations: postgres.Jsonb{},
//}
//
//func TestSetRepo_FindAllSets_whenNoRowsAreFound_returnsEmptyArray(t *testing.T) {
//	setRepo := newSetRepo(t)
//	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
//		WillReturnRows(sqlmock.NewRows(nil))
//
//	sets, err := setRepo.FindAllSets()
//	if err != nil {
//		t.Error("FindAllSets returned an error", err)
//	}
//	if len(sets) != 0 {
//		t.Error("FindAllSets returned an expected value", sets)
//	}
//}
//
//func TestSetRepo_FindAllSets_when1RowIsFound_returnsArrayWith1Value(t *testing.T) {
//	setRepo := newSetRepo(t)
//
//	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
//		WillReturnRows(sqlmock.NewRows(setColumns).
//			AddRow(basicSetValues...))
//
//	sets, err := setRepo.FindAllSets()
//	if err != nil {
//		t.Error("FindAllSets returned an error", err)
//	}
//	if len(sets) != 1 {
//		t.Error("FindAllSets returned an expected value", sets)
//	}
//	set := sets[0]
//	if !reflect.DeepEqual(*set, basicSet) {
//		t.Error("FindAllSets returned an expected value", *set, " should have been ", basicSet)
//	}
//}
//
//func TestSetRepo_FindLatestSets_whenNoRowsAreFound_returnsError(t *testing.T) {
//	setRepo := newSetRepo(t)
//	mock.ExpectQuery(regexp.QuoteMeta(sqlLatest)).
//		WillReturnRows(sqlmock.NewRows(nil))
//
//	_, err := setRepo.FindLatestSet()
//	if err == nil {
//		t.Error("FindLatestSets should have returned an error")
//	}
//}
//
//func TestSetRepo_FindLatestSets_when1RowIsFound_shouldReturnSet(t *testing.T) {
//	setRepo := newSetRepo(t)
//	mock.ExpectQuery(regexp.QuoteMeta(sqlLatest)).
//		WillReturnRows(sqlmock.NewRows(setColumns).
//			AddRow(basicSetValues...))
//
//	latestSet, err := setRepo.FindLatestSet()
//	if err != nil {
//		t.Error("FindLatestSets should have not returned an error", err)
//	}
//
//	if !reflect.DeepEqual(*latestSet, basicSet) {
//		t.Error("FindAllSets returned an expected value", *latestSet, " should have been ", basicSet)
//	}
//}
//
//func TestSetRepo_FindSet_whenNoSetIsFound_returnsError(t *testing.T) {
//	setRepo := newSetRepo(t)
//	mock.ExpectQuery(regexp.QuoteMeta(sqlFindSet)).
//		WillReturnRows(sqlmock.NewRows(nil))
//
//	_, err := setRepo.FindSet(basicSet.Code)
//	if err == nil || err.Error() != "record not found" {
//		t.Error("FindSet should have returned an error")
//	}
//}
//
//func TestSetRepo_FindSet_whenSetIsFound_shouldReturnSet(t *testing.T) {
//	setRepo := newSetRepo(t)
//	mock.ExpectQuery(regexp.QuoteMeta(sqlFindSet)).
//		WillReturnRows(sqlmock.NewRows(setColumns).
//			AddRow(basicSetValues...))
//	mock.ExpectQuery(regexp.QuoteMeta(sqlFindAllCardsInSet)).
//		WillReturnRows(sqlmock.NewRows(nil))
//	mock.ExpectQuery(regexp.QuoteMeta(sqlFindAllSheets)).
//		WillReturnRows(sqlmock.NewRows(nil))
//
//	_, err := setRepo.FindSet(basicSet.Code)
//	if err != nil {
//		t.Error("FindSet should have not returned an error", err)
//	}
//}
//
//func newSetRepo(t *testing.T) SetRepository {
//	var db *sql.DB
//	var err error
//	db, mock, err = sqlmock.New() // mock sql.DB
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	gdb, err := gorm.Open(Dialec, db) // open gorm db
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	return NewSetRepository(gdb)
//}
