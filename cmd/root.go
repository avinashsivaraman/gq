package cmd

import (
	"log"
    "fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/avinashsivaraman/gq/cmd/llm"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gq",
	Short: "A CLI to leverage Generative AI for your query",
    Long: `
  GQ (Generative Query) is a command-line tool for asking questions using AI models like Gemini.
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
    
    if len(args) == 0 && question == "" {
      return cmd.Help()  
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

    result := askQuestion(question, cmdArgs, verbose)
    write(result, os.Stdout, verbose)
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
func askQuestion(question string, data string, verbose bool) string {
    
    var extraMiddleCharacter string = "\n"
    if question == "" {
      extraMiddleCharacter = ""
    }
      
	inputQuestion := question + extraMiddleCharacter + data
    
    if verbose {
      fmt.Println("\033[33mMaking LLM Call with question: \033[0m")
      fmt.Println("\033[36m" + inputQuestion + "\033[0m")
    }

	answer, err := makeGeminiCall(inputQuestion, verbose)
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

/**
* This function writes the result as output
 */
func write(s string, w io.Writer, verbose bool) error {
	unquoted, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
      log.Fatal(err)
      return err
	}

    if verbose {
      fmt.Println("\033[32m---LLM Output---\033[0m")
    }
	_, e := fmt.Fprintln(w, unquoted)

	if e != nil {
      log.Fatal(err)
      return e
	}
	return nil
}

/**
* This function makes a call to Gemini API and retrieves the output
 */
func makeGeminiCall(question string, verbose bool) (string, error) {
	geminiLLM := llm.NewGeminiLLM()

	chatResponse, err := geminiLLM.Chat(question, verbose)
	if err != nil {
	  log.Fatal(err)
      return "", err
	}

	return chatResponse, nil
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
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/gq/.gq.yaml)")
	rootCmd.Flags().StringP("question", "q", "", "Question about the data sent")
}
