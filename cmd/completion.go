package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate completion script for your shell",
	Long: `To load completions:

Bash:

  $ source <(kwcli completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ kwcli completion bash > /etc/bash_completion.d/kwcli
  # macOS:
  $ kwcli completion bash > /usr/local/etc/bash_completion.d/kwcli

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ kwcli completion zsh > "${fpath[1]}/_kwcli"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ kwcli completion fish | source

  # To load completions for each session, execute once:
  $ kwcli completion fish > ~/.config/fish/completions/kwcli.fish

PowerShell:

  PS> kwcli completion powershell | Out-String | Invoke-Expression

  # To load completions for every session, add the output of the preceding
  # command to your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			fmt.Fprintln(os.Stdout, `# Auto-source bash-completion if available
if ! type _get_comp_words_by_ref &>/dev/null; then
	for f in /usr/share/bash-completion/bash_completion /etc/bash_completion /usr/local/etc/bash_completion; do
		if [ -f "$f" ]; then
			source "$f"
			break
		fi
	done
fi`)
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
