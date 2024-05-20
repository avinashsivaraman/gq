package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/avinashsivaraman/gq/cmd/llm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ChatProvider interface {
	Chat(string, bool) (string, error)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gq",
	Short: "A CLI to leverage Generative AI for your query",
	Long: `
  GQ (Generative AI Query) is a command-line tool for asking questions using AI models like Gemini.
  It's as easy as typing your question in the terminal or piping data for more complex queries.

  Usage examples:
    - Ask a question from the terminal:
        gq "What state is Seattle in?"

    - Process complex queries with piping:
        cat file.txt | gq -q "Explain this file to me"
    
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand(cmd, args)
	},
}

/**
* This is the main method where we control the routing of the command.
 */
func runCommand(cmd *cobra.Command, args []string) error {
	question, _ := cmd.Flags().GetString("question")
	verbose, _ := cmd.Flags().GetBool("verbose")
	provider, _ := cmd.Flags().GetString("provider")

	if len(args) == 0 && question == "" {
		var asciiArt string = `

     ░▒▓██████▓▒░   ░▒▓██████▓▒░
    ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░░▒▓█▓▒░
    ░▒▓█▓▒░        ░▒▓█▓▒░░▒▓█▓▒░
    ░▒▓█▓▒▒▓███▓▒░ ░▒▓█▓▒░░▒▓█▓▒░
    ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░░▒▓█▓▒░
    ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░░▒▓█▓▒░
     ░▒▓██████▓▒░   ░▒▓██████▓▒░
                      ░▒▓█▓▒░
                       ░▒▓██▓▒░

`
		shortDescription := `
    A CLI to leverage LLMs for your query.
    Use gq --help for more details
    `

		fmt.Println("\033[38;5;148m" + asciiArt + "\033[0m")
		fmt.Println()
		fmt.Println("\033[38;5;148m" + shortDescription + "\033[0m")
		return nil
	}

	var cmdArgs string = ""

	if isInputFromPipe() {
		if question == "" {
			return fmt.Errorf("\033[31mno question provided. Provide -q when performing Pipe operations\033[0m")
		}

		cmdArgs = readFromPipe(os.Stdin)

		if verbose {
			fmt.Println("\033[33mReading from Pipe with contents: \033[0m")
			fmt.Println("\033[36m" + cmdArgs + "\033[0m")
		}

	} else {
		if len(args) != 0 {
			cmdArgs = args[0]
		}
	}

	if provider == "" {
		provider = viper.GetString("default")
		if verbose {
			fmt.Println("\033[33mChatProvider not specified. Using default provider: \033[0m")
		}
	} else {
		if verbose {
			fmt.Println("\033[33mUsing Chat Provider: \033[0m")
		}
	}

	if verbose {
		fmt.Println("\033[36m" + provider)
		fmt.Println("\033[0m")
	}

	result := askQuestion(question, cmdArgs, provider, verbose)

	SKIP_FORMATTING := provider == "bedrock"
	write(result, os.Stdout, verbose, SKIP_FORMATTING)
	return nil
}

/*
* This function checks if the input is from pipe
 */
func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

/**
* This function asks a question to the provider and returns the answer
 */
func askQuestion(question string, data string, provider string, verbose bool) string {
	var extraMiddleCharacter string = "\n"
	if question == "" {
		extraMiddleCharacter = ""
	}

	inputQuestion := question + extraMiddleCharacter + data

	if verbose {
		fmt.Println("\033[33mMaking LLM Call with question: \033[0m")
		fmt.Println("\033[36m" + inputQuestion + "\033[0m")
	}

	chatProvider := getChatProvider(provider)

	answer, err := chatProvider.Chat(inputQuestion, verbose)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(answer, `"`)
}

/**
* This function reads the data from the pipe and asks questions and write it as output
 */
func readFromPipe(reader io.Reader) string {
	inputBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	return string(inputBytes)
}

func unquote(s string, skip bool) string {
	if skip {
		return s
	}
	unquoted, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		log.Fatal(err)
	}
	return unquoted
}

/**
* This function writes the result as output
 */
func write(s string, w io.Writer, verbose bool, skipUnquoted bool) error {
	unquoted := unquote(s, skipUnquoted)
	if verbose {
		fmt.Println("\033[32m---LLM Output---\033[0m")
	}
	_, e := fmt.Fprintln(w, unquoted)

	if e != nil {
		log.Fatal(e)
		return e
	}
	return nil
}

func getChatProvider(provider string) ChatProvider {
	switch provider {
	case "gemini":
		return llm.GeminiProvider{}
	case "openAI":
		return llm.OpenAIProvider{}
	case "azureOpenAI":
		return llm.AzureOpenAIProvider{}
	case "bedrock":
		return llm.AmznBedrockAIProvider{}
	default:
		panic("Unknown provider")
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("provider", "p", "", "the llm provider to use")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/gq/.gq.yaml)")
	rootCmd.Flags().StringP("question", "q", "", "Question about the data sent")
}
