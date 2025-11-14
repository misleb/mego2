package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/misleb/mego2/server/store"
	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := newRootCmd().ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "cmd",
		Short:         "Utility commands for the mego2 backend",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newCleanupCmd())

	return cmd
}

func newCleanupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup-data",
		Short: "Remove all tokens and users from the database (dangerous)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cleanupData(cmd.Context())
		},
	}
}

func cleanupData(ctx context.Context) error {
	db, err := store.InitDB()
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer db.Close()

	done := make(chan struct{})
	go func() {
		db.CleanupData()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	fmt.Println("Database cleanup completed.")
	return nil
}
