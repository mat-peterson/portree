package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fairy-pitta/portree/internal/config"
	"github.com/fairy-pitta/portree/internal/git"
	"github.com/fairy-pitta/portree/internal/logging"
	"github.com/spf13/cobra"
)

var (
	// Populated by PersistentPreRunE for subcommands.
	repoRoot string
	// stateRoot is the main worktree root. State (ports, PIDs, locks) must be
	// shared across all worktrees of a repo so that port allocation sees every
	// worktree's assignments; storing it per-worktree silently splits the
	// "used ports" view and allows two worktrees to claim the same port.
	stateRoot string
	cfg       *config.Config
)

// stateDir returns the shared state directory for the repository.
func stateDir() string {
	return filepath.Join(stateRoot, ".portree")
}

var rootCmd = &cobra.Command{
	Use:           "portree",
	Short:         "Git Worktree Server Manager",
	Long:          "portree manages multiple dev servers per git worktree with automatic port allocation and reverse proxy routing.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Configure log level from flags.
		verbose, _ := cmd.Flags().GetBool("verbose")
		quiet, _ := cmd.Flags().GetBool("quiet")
		if verbose {
			logging.SetLevel(logging.LevelVerbose)
		}
		if quiet {
			logging.SetLevel(logging.LevelQuiet)
		}

		// Skip repo/config detection for commands that opt out.
		if cmd.Annotations["skipRepoDetection"] == "true" {
			return nil
		}

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		repoRoot, err = git.FindRepoRoot(cwd)
		if err != nil {
			return fmt.Errorf("not inside a git repository")
		}

		logging.Verbose("repo root: %s", repoRoot)

		stateRoot, err = git.MainWorktreeRoot(cwd)
		if err != nil {
			logging.Verbose("falling back to repo root for state: %v", err)
			stateRoot = repoRoot
		}
		logging.Verbose("state root: %s", stateRoot)

		cfg, err = config.Load(repoRoot)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		logging.Verbose("loaded config with %d service(s)", len(cfg.Services))

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress all non-error output")
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "quiet")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
