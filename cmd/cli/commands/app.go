package commands

import (
	"context"

	"github.com/spf13/cobra"
	"pokergo/internal/articles"
	"pokergo/pkg/logger"
	"pokergo/pkg/timer"
)

type commandApp struct {
	*cobra.Command

	ctx    context.Context
	logger logger.Logger
	timer  timer.Timer

	artsAdapter articles.Adapter
}

func NewCommandApp(
	ctx context.Context,
	lg logger.Logger,
	tm timer.Timer,
	artsAdap articles.Adapter,
) *commandApp {
	rootCmd := &cobra.Command{
		Use:   "pokergo",
		Short: "Some useful commands for pokergo server app",
		Long: `
			The app allows to execute some useful commands,
  			like "updateArticles" etc.
			Can be easily used with k8s-jobs or any other scheduler.
			`,
	}

	app := &commandApp{
		Command:     rootCmd,
		ctx:         ctx,
		logger:      lg,
		timer:       tm,
		artsAdapter: artsAdap,
	}

	dummyCmd := &cobra.Command{
		Use:   "dummy",
		Short: "Just print debug info on the screen",
		Run: func(cmd *cobra.Command, args []string) {
			app.logger.Info("You are a dummy!")
		},
	}

	fetchArticles := &cobra.Command{
		Use:   "updateArticles",
		Short: "Updates articles in database",
		Run: func(cmd *cobra.Command, args []string) {
			app.updateArticles()
		},
	}

	rootCmd.AddCommand(dummyCmd)
	rootCmd.AddCommand(fetchArticles)
	return app
}
