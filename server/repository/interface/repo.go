package interfaces

//go:generate mockery --output=../mocks --case=underscore --all --with-expecter
type TaskManagmentInterface interface {
	TaskRepo() TaskRepo
	TaskHistoryRepo() TaskHistoryRepo
	WorkflowRepo() WorkflowRepo
	ExecutionRepo() ExecutionRepo
}
