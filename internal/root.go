/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-spec-to-scenarigo [input file]",
	Short: "API E2E Test Scenario generator",
	Long: `API E2E Test Scenario generator
This CLI tool is used to automatically generate scenarigo test formats from OpenAPI Spec. 
It generates test expectations with the results of requests to the actual API. 
			`,
	Version: version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("API E2E Test Generator")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must have API Spec")
			return
		}
		input := args[0]
		flags := *cmd.Flags()

		// 出力ファイル名を読み込み
		output, _ := flags.GetString("output-file")
		if output == "" {
			output = "scenario.yml"
		}

		// テストケースの指定を読み込み
		casesStr, _ := flags.GetString("test-case")
		var cases []string
		if casesStr != "" {
			cases = strings.Split(casesStr, ",")
		} else {
			cases = []string{}
		}
		// API Specを読み込み
		result, err := GenItem(input, cases)
		if err != nil {
			fmt.Println("error")
		}

		// APIエンドポイントを読み込み
		host, _ := flags.GetString("host")
		if host != "" {
			result.BaseUrl = host
		}

		// 追加のテストパターンを読み込む
		paramfile, _ := flags.GetString("csv-file")
		var param *map[string]addParam
		if paramfile != "" {
			param, err = AddParam(paramfile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// dry run flagの読み込み
		dryrunFlg, _ := flags.GetBool("dry-run")

		// シナリオを作成
		if dryrunFlg {
			fmt.Println(input)
			fmt.Println(output)
			fmt.Println("Title : ", result.Title)
			fmt.Println("Description : ", result.Description)
			fmt.Println("Version : ", result.Version)
			fmt.Println("BaseURL : ", result.BaseUrl)

			paths := result.PathSpec
			fmt.Println("==========================")
			for _, r := range paths {
				fmt.Println(r.Path)
				for _, m := range r.Methods {
					fmt.Println("Method : ", m.Method)
					fmt.Println("Request Body : ", m.Body)
					fmt.Println("Request Param : ", m.Params)
					fmt.Println("Response : ", m.Response)
					fmt.Println("==========================")
				}
			}
		} else if param != nil {
			err = GenScenario(result, output, param)
			if err != nil {
				fmt.Println("error:", err)
			}
		} else {
			err = GenScenario(result, output)
			if err != nil {
				fmt.Println("error:", err)
			}
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("done")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	version = v
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("help", "h", false, "Help message")
	rootCmd.Flags().StringP("output-file", "o", "", "output file name")
	rootCmd.Flags().StringP("host", "s", "", "API EndPoint")
	rootCmd.Flags().BoolP("dry-run", "d", false, "dry run mode. not generate scenario file")
	rootCmd.Flags().StringP("csv-file", "c", "", "add test pattern parameter")
	rootCmd.Flags().StringP("test-case", "t", "", "determine use test case")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
