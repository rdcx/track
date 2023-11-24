/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"track/types"

	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Please provide a domain and key")
			return
		}

		domain := args[0]
		key := args[1]

		res, err := http.Get("https://trackcmd.com/hits/" + domain + "?k=" + key)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		type hitsResponse struct {
			Hits types.UrlMap `json:"hits"`
		}

		var hr hitsResponse

		err = json.Unmarshal(data, &hr)
		if err != nil {
			fmt.Println(err)
			return
		}

		// group by date
		// past 24 hours
		past24Hours := time.Now().AddDate(0, 0, -1)
		// past 7 days
		past7Days := time.Now().AddDate(0, 0, -7)
		// past 30 days
		past30Days := time.Now().AddDate(0, 0, -30)
		// past 90 days
		past90Days := time.Now().AddDate(0, 0, -90)
		days := make(map[string]map[string]map[string]int64)

		for url, v := range hr.Hits {
			for _, h := range v {

				if h.Time.After(past24Hours) {
					if _, ok := days["past24Hours"]; !ok {
						days["past24Hours"] = make(map[string]map[string]int64)
					}

					if _, ok := days["past24Hours"][string(url)]; !ok {
						days["past24Hours"][string(url)] = make(map[string]int64)
					}

					if _, ok := days["past24Hours"][string(url)]["total"]; !ok {
						days["past24Hours"][string(url)]["total"] = 0
					}

					days["past24Hours"][string(url)][h.Loc]++
					days["past24Hours"][string(url)]["total"]++
				}

				if h.Time.After(past7Days) {
					if _, ok := days["past7Days"]; !ok {
						days["past7Days"] = make(map[string]map[string]int64)
					}
					if _, ok := days["past7Days"][string(url)]; !ok {
						days["past7Days"][string(url)] = make(map[string]int64)
					}

					if _, ok := days["past7Days"][string(url)]["total"]; !ok {
						days["past7Days"][string(url)]["total"] = 0
					}
					days["past7Days"][string(url)][h.Loc]++
					days["past7Days"][string(url)]["total"]++
				}

				if h.Time.After(past30Days) {
					if _, ok := days["past30Days"]; !ok {
						days["past30Days"] = make(map[string]map[string]int64)
					}

					if _, ok := days["past30Days"][string(url)]; !ok {
						days["past30Days"][string(url)] = make(map[string]int64)
					}

					if _, ok := days["past30Days"][string(url)]["total"]; !ok {
						days["past30Days"][string(url)]["total"] = 0
					}

					days["past30Days"][string(url)][h.Loc]++
					days["past30Days"][string(url)]["total"]++
				}

				if h.Time.After(past90Days) {
					if _, ok := days["past90Days"]; !ok {
						days["past90Days"] = make(map[string]map[string]int64)
					}

					if _, ok := days["past90Days"][string(url)]; !ok {
						days["past90Days"][string(url)] = make(map[string]int64)
					}

					if _, ok := days["past90Days"][string(url)]["total"]; !ok {
						days["past90Days"][string(url)]["total"] = 0
					}
					days["past90Days"][string(url)][h.Loc]++
					days["past90Days"][string(url)]["total"]++
				}
			}
		}

		pp.Println(days)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
