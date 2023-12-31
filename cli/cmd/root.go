/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"track/types"

	"github.com/spf13/cobra"
)

var checkEmoji = "✅"
var crossEmoji = "❌"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "track",
	Short: "Start tracking traffic to your websites.",
	Long: `You can track your websites, all you need to do provide us with your domain name. It will return a JS snippet that you can add to your website. 
For example:
	track yourdomain.com

<script>
	(function() {
		var url = window.location.href;
		fetch("https://trackcmd.com/hit/" + btoa(url));
	})();
</script>
		`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var domain string
		if len(args) == 0 {
			// TODO: prompt user for domain
			fmt.Println("Please provide the website domain:")
			fmt.Scanln(&domain)

			if strings.Contains(domain, "http") {
				fmt.Println("Please provide the domain name without the protocol (http:// or https://) for example: yourwebsite.com")
				os.Exit(1)
			}

			domain = strings.TrimSpace(domain)
		} else {
			domain = args[0]
		}

		res, err := http.Post("https://trackcmd.com/track/"+domain, "text/plain", nil)

		if err != nil {
			panic(err)
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {

			body, err := io.ReadAll(res.Body)

			if err != nil {
				panic(err)
			}

			var msg types.MessageResponse

			err = json.Unmarshal(body, &msg)

			if err != nil {
				panic(err)
			}

			fmt.Println(crossEmoji + " " + strings.Title(msg.Message))
			os.Exit(1)
		}

		var resp types.TrackResponse

		body, err := io.ReadAll(res.Body)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &resp)

		if err != nil || !resp.Success || resp.Message != "ok" || resp.Key == "" {
			panic("failed to track domain")
		}

		res, err = http.Get("http://trackcmd.com/tracker?k=" + string(resp.Key))

		if err != nil {
			panic(err)
		}

		defer res.Body.Close()

		fmt.Printf(checkEmoji+" Now tracking %s, use the following snippet on all pages:\n", domain)

		body, err = io.ReadAll(res.Body)

		if err != nil {
			panic(err)
		}

		fmt.Println(string(body))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
