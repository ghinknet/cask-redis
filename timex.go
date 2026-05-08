package redis

import (
	"time"

	"go.gh.ink/timex"
)

const minExpireDuration = 1 * time.Millisecond

func ToSetPXDuration(d timex.Duration) time.Duration {
	std, inf := d.ToStdDuration()

	if inf == timex.PosInfTime {
		return 0
	}

	if inf == timex.FiniteTime && std > 0 {
		return std
	}

	// Expire immediately
	return minExpireDuration
}

func FromTTLDuration(ttl time.Duration) timex.Duration {
	switch ttl {
	case -2:
		return timex.NewNegInfDuration()
	case -1:
		return timex.NewPosInfDuration()
	default:
		return timex.FromStdDuration(ttl)
	}
}
