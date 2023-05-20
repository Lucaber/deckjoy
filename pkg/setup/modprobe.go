package setup

import (
	"context"
	"github.com/lucaber/deckjoy/pkg/util"
)

func Modprobe(module string) error {
	return util.Exec(context.Background(), "modprobe", module)
}
