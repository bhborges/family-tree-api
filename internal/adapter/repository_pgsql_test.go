package adapter

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bhborges/family-tree-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostgresRepository_BuildFamilyTree(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresRepository(sqlDB, nil)

	parent := &domain.Person{
		Name: "Parent",
	}
	child := &domain.Person{
		Name: "Child",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"people\"").
		WithArgs(parent.ID, parent.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO \"people\"").
		WithArgs(child.ID, child.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO \"relationships\"").
		WithArgs(parent.ID, child.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// if _, err := repo.CreateRelationship(context.Background(), domain.Relationship{
	// 	ParentID: parent.ID,
	// 	ChildID:  child.ID,
	// }); err != nil {
	// 	t.Fatal(err)
	// }

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(parent.ID, parent.Name)

	mock.ExpectQuery("SELECT (.+) FROM \"people\" WHERE id = (.+)").
		WithArgs(parent.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "name"}).AddRow(child.ID, child.Name)

	mock.ExpectQuery("SELECT (.+) FROM \"people\" WHERE id = (.+)").
		WithArgs(child.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "parent_id", "child_id"}).
		AddRow(1, parent.ID, child.ID)

	mock.ExpectQuery("SELECT (.+) FROM \"relationships\" WHERE parent_id = (.+)").
		WithArgs(parent.ID).
		WillReturnRows(rows)

	tree, err := repo.BuildFamilyTree(context.Background(), parent.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(tree.Members))
	assert.Equal(t, parent.Name, tree.Members[0].Name)
	assert.Equal(t, 1, len(tree.Members[0].Relationships))
	assert.Equal(t, child.Name, tree.Members[0].Relationships[0].Name)
}
