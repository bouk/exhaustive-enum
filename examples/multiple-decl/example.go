// In this example the exhaustive enum tag should only apply to the Day type. The first exhaustive-enum comment doesn't do anything
package example

//exhaustive-enum
type (
	//exhaustive-enum
	Day  int
	Time int
)

const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

const (
	A Time = iota
	B
	C
	D
	E
)

func a(day Day) {
	switch day {
	case Monday:
	case Tuesday:
	case Wednesday:
	case Thursday:
	case Friday:
	case Saturday:
		//case Sunday:
	}
}
func b(t Time) {
	switch t {
	case A:
	case B:
	case C:
	case D:
		//case E:
	}
}
