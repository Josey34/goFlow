package safesync

// UNSAFE - This WILL have a race condition!
// DO NOT use this in production - only for learning
type UnsafeMap map[string]int

func RaceExample() {
	unsafe := make(UnsafeMap)

	// Goroutine 1: writes
	go func() {
		for i := 0; i < 1000; i++ {
			unsafe["key"] = i // WRITE - no lock!
		}
	}()

	// Goroutine 2: reads
	go func() {
		for i := 0; i < 1000; i++ {
			_ = unsafe["key"] // READ - no lock!
		}
	}()
}
