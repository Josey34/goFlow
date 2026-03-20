package cmd

import (
	"context"
	"fmt"
	"log"

	"goflow/factory"
	"goflow/pipeline"
	"goflow/worker"
)

func Process(ctx context.Context, f *factory.Factory) error {
	fmt.Println("Starting document processing pipeline...")
	fmt.Printf("Workers: %d | ChunkSize: %d | ChunkOverlap: %d\n",
		f.Config.Workers, f.Config.ChunkSize, f.Config.ChunkOverlap)

	pool := worker.NewPool(f.Config.Workers, f.SQSConsumer, f.MinIODownloader, f.Config.MaxRetries, f.Limiter)

	extractor := pipeline.NewExtractor()
	chunker := pipeline.NewChunker(pipeline.StageConfig{
		ChunkSize:    f.Config.ChunkSize,
		ChunkOverlap: f.Config.ChunkOverlap,
	})
	deduplicator := pipeline.NewDeduplicator(f.ResultRepository, f.Cache)
	aggregator := pipeline.NewAggregator()
	writer := pipeline.NewWriter(f.ResultRepository, f.ChunkRepository)

	pool.Start()

	successCount := 0
	errorCount := 0

	go func() {
		for {
			select {
			case <-ctx.Done():
				pool.Stop()
				return
			case result := <-pool.Results():
				if result == nil {
					return
				}

				var stageInput interface{} = result
				stages := []pipeline.Stage{extractor, chunker, deduplicator, aggregator, writer}

				for _, stage := range stages {
					output, err := stage.Process(ctx, stageInput)
					if err != nil {
						log.Printf("❌ %s error: %v", stage.Name(), err)
						errorCount++
						break
					}
					stageInput = output
				}

				successCount++
				fmt.Printf("✓ Processed document %s\n", result.Task.Event.DocumentID)

			case err := <-pool.Errors():
				if err != nil {
					log.Printf("❌ Pool error: %v", err)
					errorCount++
				}
			}
		}
	}()

	<-ctx.Done()
	fmt.Printf("\nShutdown - Processed: %d | Errors: %d\n", successCount, errorCount)
	return nil
}
