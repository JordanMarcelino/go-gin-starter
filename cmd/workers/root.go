package workers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/logger"
	"github.com/JordanMarcelino/go-gin-starter/internal/provider"
	"github.com/spf13/cobra"
)

func Start() {
	cfg := config.InitConfig()
	logger.SetZerologLogger(cfg)
	provider.InitGlobal(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd := &cobra.Command{}
	cmd := []*cobra.Command{
		{
			Use:   "serve-all",
			Short: "Run all",
			Run: func(cmd *cobra.Command, _ []string) {
				runHttpWorker(cfg, ctx)
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
