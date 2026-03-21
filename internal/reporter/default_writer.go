package reporter

import (
	"fmt"
	"io"

	"github.com/velosypedno/resource-allocation/internal/scheduler"
)

type Formatter interface {
	Format(results []scheduler.PlanResult) (string, error)
}

type Reporter struct {
	writer    io.Writer
	formatter Formatter
}

func NewReporter(w io.Writer, f Formatter) *Reporter {
	return &Reporter{
		writer:    w,
		formatter: f,
	}
}

func (r *Reporter) Generate(results []scheduler.PlanResult) error {
	content, err := r.formatter.Format(results)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(r.writer, content)
	return err
}
