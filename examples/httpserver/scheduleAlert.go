package poker

import (
	"fmt"
	"time"
)

type ScheduledAlert struct {
	At time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}
