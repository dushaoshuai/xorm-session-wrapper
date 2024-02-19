package wrapper

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"slices"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

const (
	testDBName  = "test_db"
	testAddress = "localhost"
	testPort    = 13306
)

type testTable struct {
	ID int64 `xorm:"'id' not null pk autoincr comment('primary key') BIGINT"`
}

func (m *testTable) TableName() string {
	return "test_table"
}

var (
	test1    = testTable{ID: 1}
	test2    = testTable{ID: 2}
	test3    = testTable{ID: 3}
	test4    = testTable{ID: 4}
	test5    = testTable{ID: 5}
	test6    = testTable{ID: 6}
	test7    = testTable{ID: 7}
	allTests = []testTable{test1, test2, test3, test4, test5, test6, test7}
)

func mustGetXormEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql",
		fmt.Sprintf("tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_general_ci",
			testAddress, testPort, testDBName),
	)
	if err != nil {
		panic(err)
	}

	return engine
}

func createTestServer() (run, clos func()) {
	memoryDB := memory.NewDatabase(testDBName)
	memoryDB.EnablePrimaryKeyIndexes()

	engine := sqle.NewDefault(memory.NewDBProvider(memoryDB))

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", testAddress, testPort),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	return func() {
			if err = s.Start(); err != nil {
				panic(err)
			}
		}, func() {
			s.Close()
		}
}

func createTestTable() {
	engine := mustGetXormEngine()
	err := engine.CreateTables(testTable{})
	if err != nil {
		panic(err)
	}

	affected, err := engine.Insert(&allTests)
	if err != nil {
		panic(err)
	}
	if affected != int64(len(allTests)) {
		panic("wrong inserted records number")
	}
}

func TestMain(m *testing.M) {
	runServer, closeServer := createTestServer()
	defer closeServer()
	go runServer()
	time.Sleep(2 * time.Second)

	createTestTable()

	os.Exit(m.Run())
}

func TestSession_In(t *testing.T) {
	engine := mustGetXormEngine()
	sess := NewSession(engine.NewSession())

	tests := []struct {
		name   string
		values []any
		want   []testTable
	}{
		{currLine(), nil, allTests},
		{currLine(), []any{}, allTests},
		{currLine(), []any{nil}, allTests},
		{currLine(), []any{nil, nil}, allTests},
		{currLine(), []any{[]string(nil)}, allTests},
		{currLine(), []any{[]string{}}, allTests},
		{currLine(), []any{[]string{"a"}}, nil},
		{currLine(), []any{[]int64(nil)}, allTests},
		{currLine(), []any{[]int64{}}, allTests},
		{currLine(), []any{[]int64{1}}, []testTable{test1}},
		{currLine(), []any{[]int64{1, 2}}, []testTable{test1, test2}},
		{currLine(), []any{[]int64{1, 2, 3}}, []testTable{test1, test2, test3}},
		{currLine(), []any{[]int64{1, 2, 3, 4}}, []testTable{test1, test2, test3, test4}},
		{currLine(), []any{[]int64{1, 2, 3, 4, 5, 6, 7}}, allTests},
		{currLine(), []any{[]int64{1, 2, 3, 4, 5, 6, 7, 8}}, allTests},
		{currLine(), []any{[]int64{1, 77}}, []testTable{test1}},
		{currLine(), []any{[]int64{77}}, nil},
		{currLine(), []any{1}, []testTable{test1}},
		{currLine(), []any{1, 2}, []testTable{test1, test2}},
		{currLine(), []any{1, 2, 3}, []testTable{test1, test2, test3}},
		{currLine(), []any{1, 2, 3, 4}, []testTable{test1, test2, test3, test4}},
		{currLine(), []any{1, 2, 3, 4, 5, 6, 7}, allTests},
		{currLine(), []any{1, 2, 3, 4, 5, 6, 7, 8}, allTests},
		{currLine(), []any{1, "a"}, []testTable{test1}},
		{currLine(), []any{1, 77}, []testTable{test1}},
		{currLine(), []any{77}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testTable
			err := sess.In("id", tt.values...).Find(&got)
			if err != nil {
				t.Fatalf("In() error = %v", err)
			}
			if !deepEqual(got, tt.want) {
				t.Errorf("In() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Equal(t *testing.T) {
	engine := mustGetXormEngine()
	sess := NewSession(engine.NewSession())

	tests := []struct {
		name  string
		value any
		want  []testTable
	}{
		{currLine(), nil, allTests},
		{currLine(), "", allTests},
		{currLine(), "a", nil},
		{currLine(), -1, nil},
		{currLine(), 0, allTests},
		{currLine(), 1, []testTable{test1}},
		{currLine(), 2, []testTable{test2}},
		{currLine(), 77, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testTable
			err := sess.Equal("id", tt.value).Find(&got)
			if err != nil {
				t.Fatalf("Equal() error = %v", err)
			}
			if !deepEqual(got, tt.want) {
				t.Errorf("Equal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Between(t *testing.T) {
	engine := mustGetXormEngine()
	sess := NewSession(engine.NewSession())

	tests := []struct {
		name   string
		ranger *Ranger
		want   []testTable
	}{
		{currLine(), nil, allTests},
		{currLine(), &Ranger{1, 2}, []testTable{test1, test2}},
		{currLine(), &Ranger{1, 7}, allTests},
		{currLine(), &Ranger{-1, 8}, allTests},
		{currLine(), &Ranger{2, 1}, nil},
		{currLine(), &Ranger{8, -1}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testTable
			err := sess.Between("id", tt.ranger).Find(&got)
			if err != nil {
				t.Fatalf("Between() error = %v", err)
			}
			if !deepEqual(got, tt.want) {
				t.Errorf("Between() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func deepEqual(x, y []testTable) bool {
	cmp := func(a, b testTable) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	}
	slices.SortFunc(x, cmp)
	slices.SortFunc(y, cmp)

	return reflect.DeepEqual(x, y)
}

func currLine() string {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("line%d", line)
}
