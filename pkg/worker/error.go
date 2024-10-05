package worker

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

type CustomErrorHandler struct{}

func (*CustomErrorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
	fmt.Printf("Job errored with: %s\n", err)
	return nil
}

func (*CustomErrorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any, trace string) *river.ErrorHandlerResult {
	fmt.Printf("Job panicked with: %v\n", panicVal)
	fmt.Printf("Job panicked with: %v\n", trace)

	// Either function can also set the job to be immediately cancelled.
	return &river.ErrorHandlerResult{SetCancelled: true}
}
