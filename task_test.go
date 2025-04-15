package go_ticks

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/parametalol/go-ticks/ticker"
	"github.com/parametalol/go-ticks/utils"
)

func TestTask(t *testing.T) {
	t.Run("failed function", func(t *testing.T) {
		ticker := ticker.NewTimer(time.Second)
		defer ticker.Stop()

		NewTask(ticker, utils.WithLog[time.Time](os.Stdout, os.Stdout, "test", func() error {
			fmt.Println("tick")
			return utils.ErrStopped
		}))
		// Output:
		// tick
	})
}
