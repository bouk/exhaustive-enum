package example2

import (
	"github.com/bouk/exhaustive-enum/examples/example"
)

func a(day example.Day) {
	switch day {
	case example.Monday:
	case example.Tuesday:
	case example.Wednesday:
	case example.Thursday:
	case example.Friday:
	case example.Saturday:
	case example.Sunday:
	}
}
