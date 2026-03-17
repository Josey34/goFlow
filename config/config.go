package config

import "flag"

type Config struct {
	DBPath string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool

	SQSEndpoint string
	SQSRegion   string
	SQSQueueUrl string

	Workers      int
	ChunkSize    int
	ChunkOverlap int
	MaxRetries   int

	CacheTTL  int
	RateLimit int

	Verbose bool
}

func Parse() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.DBPath, "db", "./docvault.db", "Path to shared SQLite DB")

	flag.StringVar(&cfg.MinIOEndpoint, "minio-endpoint", "localhost:9000", "MinIO endpoint")
	flag.StringVar(&cfg.MinIOAccessKey, "minio-access-key", "minioadmin", "MinIO access key")
	flag.StringVar(&cfg.MinIOSecretKey, "minio-secret-key", "minioadmin", "MinIO secret key")
	flag.StringVar(&cfg.MinIOBucket, "minio-bucket", "docvault", "MinIO bucket name")
	flag.BoolVar(&cfg.MinIOUseSSL, "minio-ssl", false, "Use HTTPS for MinIO")

	flag.StringVar(&cfg.SQSEndpoint, "sqs-endpoint", "http://localhost:4566", "SQS endpoint (optional for AWS)")
	flag.StringVar(&cfg.SQSRegion, "sqs-region", "us-east-1", "AWS region")
	flag.StringVar(&cfg.SQSQueueUrl, "sqs-queue-url", "http://localhost:4566/000000000000/docvault-events", "SQS queue URL")

	flag.IntVar(&cfg.Workers, "workers", 3, "Number of concurrent workers")
	flag.IntVar(&cfg.ChunkSize, "chunk-size", 1000, "Characters per chunk")
	flag.IntVar(&cfg.ChunkOverlap, "chunk-overlap", 200, "Overlap between chunks")
	flag.IntVar(&cfg.MaxRetries, "max-retries", 3, "Max retries on failure")

	flag.IntVar(&cfg.CacheTTL, "cache-ttl", 3600, "Cache TTL in seconds")
	flag.IntVar(&cfg.RateLimit, "rate-limit", 5, "Max concurrent downloads")

	flag.BoolVar(&cfg.Verbose, "verbose", false, "Verbose logging")

	flag.Parse()
	return cfg
}
