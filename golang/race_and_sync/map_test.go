package race_and_sync

import "testing"

func TestMapIsNonConcurrencySafe(t *testing.T) {
	MapIsNonConcurrencySafe()
}

func TestConcurrencySafeMapWithLock(t *testing.T) {
	ConcurrencySafeMapWithLock()
}

func TestSyncMapIsConcurrencySafe(t *testing.T) {
	SyncMapIsConcurrencySafe()
}
