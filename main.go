package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/0xTatsu/goco/pkg"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("cannot create gcs client: %s", err)
	}

	gcsPkg := pkg.NewGCS(
		storageClient,
		viper.GetString("BUCKET_NAME"),
	)

	fmt.Println(gcsPkg.IsExistent(ctx, "dir/filename.csv"))
}
