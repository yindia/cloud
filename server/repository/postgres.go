package repositories

import (
	"fmt"
	gormimpl "task/server/repository/gormimpl"
	interfaces "task/server/repository/interface"

	"gorm.io/gorm"
)

type Postgres struct {
	task      interfaces.TaskRepo
	history   interfaces.TaskHistoryRepo
	workflow  interfaces.WorkflowRepo
	execution interfaces.ExecutionRepo
}

func (r Postgres) TaskRepo() interfaces.TaskRepo {
	fmt.Println(r.task)
	return r.task
}

func (r Postgres) TaskHistoryRepo() interfaces.TaskHistoryRepo {
	return r.history
}

func (r Postgres) WorkflowRepo() interfaces.WorkflowRepo {
	return r.workflow
}

func (r Postgres) ExecutionRepo() interfaces.ExecutionRepo {
	return r.execution
}

func NewPostgresRepo(db *gorm.DB) interfaces.TaskManagmentInterface {
	return &Postgres{
		task:      gormimpl.NewTaskRepo(db),
		history:   gormimpl.NewTaskHistoryRepo(db),
		workflow:  gormimpl.NewWorkflowRepo(db),
		execution: gormimpl.NewExecutionRepo(db),
	}
}
