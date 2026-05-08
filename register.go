package redis

import (
	"go.gh.ink/cask/driver"
)

func init() {
	driver.Register(Name, Driver{})
}
