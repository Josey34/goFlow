package cmd

import (
	"context"
	"fmt"
	"goflow/factory"
)

func Health(ctx context.Context, f *factory.Factory) error {
	fmt.Println("Checking system health...")

	if err := f.DB.PingContext(ctx); err != nil {
		fmt.Println("SQLite: Failed -", err)
		return err
	}

	fmt.Println("SQLite: Connected")

	if f.Config.SQSQueueUrl == "" {
		fmt.Println("SQS: Queue URL not configured")
		return fmt.Errorf("SQS_QUEUE_URL not set")
	}

	fmt.Println("SQS: Configured - ", f.Config.SQSQueueUrl)

	fmt.Println("\n All system operational")
	return nil
}
