package repositories

import (
	"database/sql"
	"fmt"
	gormimpl "task/server/repository/gormimpl"
	interfaces "task/server/repository/interface"

	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type Postgres struct {
	task    interfaces.TaskRepo
	history interfaces.TaskHistoryRepo
}

func (r Postgres) TaskRepo() interfaces.TaskRepo {
	fmt.Println(r.task)
	return r.task
}

func (r Postgres) TaskHistoryRepo() interfaces.TaskHistoryRepo {
	return r.history
}

func NewPostgresRepo(db *gorm.DB, riverClient *river.Client[*sql.Tx]) interfaces.TaskManagmentInterface {
	return &Postgres{
		task:    gormimpl.NewTaskRepo(db, riverClient),
		history: gormimpl.NewTaskHistoryRepo(db, riverClient),
	}
}
