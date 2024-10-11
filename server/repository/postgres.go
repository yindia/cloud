package repositories

import (
	"fmt"
	gormimpl "task/server/repository/gormimpl"
	interfaces "task/server/repository/interface"

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

func NewPostgresRepo(db *gorm.DB) interfaces.TaskManagmentInterface {
	return &Postgres{
		task:    gormimpl.NewTaskRepo(db),
		history: gormimpl.NewTaskHistoryRepo(db),
	}
}
