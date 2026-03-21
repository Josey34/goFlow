package cmd

import (
	"context"
	"fmt"

	"goflow/factory"
)

func Stats(f *factory.Factory) error {
	ctx := context.Background()

	stats, err := f.ResultRepository.GetStats(ctx)
	if err != nil {
		return err
	}

	cacheStats := f.Cache.Stats()

	fmt.Println("\n=== Processing Statistics ===")
	fmt.Printf("Total Processed:  %d\n", stats.TotalProcessed)
	fmt.Printf("Duplicates Found: %d\n", stats.DuplicatesFound)
	fmt.Printf("Errors:           %d\n", stats.ErrorsEncountered)
	fmt.Printf("Avg Duration:     %.2f seconds\n", stats.AvgProcessingTime)

	fmt.Println("\n=== Cache Statistics ===")
	fmt.Printf("Cached Items: %d\n", cacheStats.Items)
	fmt.Printf("Cache Hits:   %d\n", cacheStats.Hits)
	fmt.Printf("Cache Misses: %d\n", cacheStats.Misses)
	if cacheStats.Hits+cacheStats.Misses > 0 {
		hitRate := float64(cacheStats.Hits) / float64(cacheStats.Hits+cacheStats.Misses) * 100
		fmt.Printf("Hit Rate:     %.2f%%\n", hitRate)
	}

	fmt.Println()
	return nil
}
