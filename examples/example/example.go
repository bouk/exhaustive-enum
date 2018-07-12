package example

//exhaustive-enum
type Day int

const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
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
