package graph

import (
	"time"

	"github.com/artemmarkaryan/exlex-backend/graph/model"
)

func timeToDate(t time.Time) (d model.Date) {
	d.Year = t.Year()
	d.Month = int(t.Month())
	d.Day = t.Day()
	return
}
