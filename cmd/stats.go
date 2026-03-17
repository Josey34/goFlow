package cmd

import (
	"fmt"
	"goflow/factory"
)

func Stats(f *factory.Factory) error {
	fmt.Println("=== Cache Statistics ===")

	stats := f.Cache.Stats()

	fmt.Printf("Cache Hits:    %d\n", stats.Hits)
	fmt.Printf("Cache Misses:  %d\n", stats.Misses)
	fmt.Printf("Items Cached:  %d\n", stats.Items)

	total := stats.Hits + stats.Misses
	if total > 0 {
		hitRate := float64(stats.Hits) / float64(total) * 100
		fmt.Printf("Hit Rate:      %.2f%%\n", hitRate)
	} else {
		fmt.Println("Hit Rate:      N/A (no accesses)")
	}

	return nil
}
