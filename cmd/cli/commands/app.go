package commands

import (
	"github.com/spf13/cobra"
	"pokergo/internal/articles"
	"pokergo/internal/mongo"
	"pokergo/pkg/logger"
	"pokergo/pkg/timer"
)

type commandApp struct {
	*cobra.Command

	logger logger.Logger
	timer  timer.Timer

	mongoColls  *mongo.Collections
	artsAdapter articles.Adapter
}

func NewCommandApp(
	lg logger.Logger,
	tm timer.Timer,
	artsAdap articles.Adapter,
	mongoColls *mongo.Collections,
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
		logger:      lg,
		timer:       tm,
		artsAdapter: artsAdap,
		mongoColls:  mongoColls,
	}

	dummyCmd := &cobra.Command{
		Use:   "dummy",
		Short: "Just print debug info on the screen",
		Run: func(cmd *cobra.Command, args []string) {
			app.logger.Info("You are a dummy!")
		},
	}

	mongoIndexes := &cobra.Command{
		Use:   "mongoIndexes",
		Short: "sets mongo indexes (must be run after each index change)",
		Run: func(cmd *cobra.Command, args []string) {
			app.mongoIndexes()
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
	rootCmd.AddCommand(mongoIndexes)
	rootCmd.AddCommand(fetchArticles)

	return app
}
