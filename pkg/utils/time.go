package utils

import "time"

var (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02T15:04:05"
)

var jakartaLoc *time.Location

func init() {
	var err error
	jakartaLoc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		jakartaLoc = time.FixedZone("Asia/Jakarta", 7*3600)
	}
}

// ====== Parsing ======
func ParseDateWIB(dateStr string) (time.Time, error) {
	return time.ParseInLocation(DateLayout, dateStr, jakartaLoc)
}

func ParseDateTimeWIB(dateTimeStr string) (time.Time, error) {
	return time.ParseInLocation(DateTimeLayout, dateTimeStr, jakartaLoc)
}

// ====== View/Output ======
func FormatDateForView(t time.Time, layout string) string {
	tz, _ := TimeIn(t, "Asia/Jakarta")
	return tz.Format(layout)
}

func TimeIn(t time.Time, name string) (time.Time, error) {
	location, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(location)
	}
	return t, err
}

func JakartaLocation() *time.Location {
	return jakartaLoc
}
