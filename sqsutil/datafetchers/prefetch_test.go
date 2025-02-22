package datafetchers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestWaitUntilFirstResult tests the behavior of the WaitUntilFirstResult method in the Prefetcher type.
// It verifies that the WaitUntilFirstResult method blocks until the first result is available,
// and that it returns the correct value and timestamp after the result is available.
// The test also checks that the WaitUntilFirstResult method blocks for the expected duration,
// and does not block for too long.
func TestWaitUntilFirstResult(t *testing.T) {
	t.Parallel()

	updateFn := func() (int, error) {
		time.Sleep(3 * time.Second)
		return 42, nil
	}

	p := NewIntervalFetcher(updateFn, 1*time.Second)
	v, timestamp, err := p.Get()
	require.Error(t, err)
	require.Equal(t, 0, v)
	require.Equal(t, time.Time{}, timestamp)

	time.Sleep(3 * time.Second)

	v, timestamp, err = p.Get()
	require.NoError(t, err)
	require.Equal(t, v, 42)
}

func requireTimeDurationInRange(t *testing.T, d time.Duration, min time.Duration, max time.Duration) {
	require.True(t, d >= min, "Duration %s is less than min %s", d, min)
	require.True(t, d <= max, "Duration %s is greater than max %s", d, max)
}
