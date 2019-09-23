package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var plumberEntryPoint string

func init() {
	rootCmd.AddCommand(routesCmd)
	routesCmd.Flags().StringVarP(&plumberEntryPoint, "entry", "e", "entrypoint.r", "Plumber application entrypoint file")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Display all routes in your Plumber application",
	Long:  `A quick way to visualize your application's routing structure`,
	Run: func(cmd *cobra.Command, args []string) {
		RouteStructure(plumberEntryPoint, hostValue, portValue, absoluteHost, routeFilter)
	},
}

// RouteStructure outputs the parsed endpoints for a given entrypoint file
func RouteStructure(plumberEntryPoint string, host string, port int, absoluteHost bool, routeFilter string) {

	var dirPath string

	if watchDir != "" {
		dirPath = filepath.Base(watchDir)
	} else {
		// watch current
		cwd, _ := os.Getwd()
		dirPath = filepath.Base(cwd)

	}

	// gen route structure, maybe write a lexer in the future
	plumb, _ := regexp.Compile(`(?i)(?P<comment>#*).*plumb\("(?P<plumber>[a-zA-Z0-9_]+\.[rR])"\)`)

	routes, _ := regexp.Compile(`(?i)#\*\s*@(get|post|put|delete|head)\s/[a-zA-Z0-9\-_\/<>:]+`)
	assets, _ := regexp.Compile(`(?i)#\*\s*@assets\s*[\.\/a-zA-Z0-9\_]+\s[\.\/a-zA-Z0-9\_]*`)

	// other components
	programmaticRoutes, _ := regexp.Compile(`(?i)\$handle\(\"(get|post|put|delete|head)\",\s*\"\/(?P<route>[a-zA-Z0-9_]+)\"`)
	// mountedRoutes, _ := regexp.Compile(`#\* @(get|post|put|delete|head)\s/[a-zA-Z0-9\-_\/<>:]+`)
	// mountedAssets, _ := regexp.Compile(`#\* @assets\s[\.\/a-zA-Z0-9\_]+\s[\.\/a-zA-Z0-9\_]*`)

	dat, err := ioutil.ReadFile(plumberEntryPoint)
	check(err)

	// if length > 0 then try to read the routes file
	// figure out nests and mounts
	entryMatches := plumb.FindAllStringSubmatch(string(dat), -1)
	// index remains the same if no match
	// loop through entry matches
	if len(entryMatches) > 0 {
		for _, entry := range entryMatches {

			comment := entry[1]
			if comment != "#" {

				dat, err := ioutil.ReadFile(entry[2])
				check(err)

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

				// route table
				// refactor into function
				fmt.Printf("[%s] Routes: \n\n", dirPath)
				routeMatches := routes.FindAllStringSubmatch(string(dat), -1)
				for _, match := range routeMatches {
					s := strings.TrimPrefix(match[0], "#*")
					parts := strings.Split(s, " ")

					// route filter
					var printRoute = true
					if routeFilter != "" && !strings.Contains(parts[2], routeFilter) {
						printRoute = false
					}

					// flag for absolute endpoint
					// needs refactored into function
					if printRoute {
						if absoluteHost {
							var endpoint string
							if host != "" {
								endpoint = strings.TrimRight(host, "/") + ":" + strconv.Itoa(port) + parts[2]
							} else {
								endpoint = parts[2]
							}

							fmt.Fprintln(w, strings.ToUpper(parts[1])+"\t"+endpoint)
						} else {
							fmt.Fprintln(w, strings.Join(parts, "\t"))
						}
					}

				}

				// programmatic routes
				programmaticRouteMatches := programmaticRoutes.FindAllStringSubmatch(string(dat), -1)
				for _, match := range programmaticRouteMatches {
					s := strings.TrimPrefix(match[0], "$handle(")
					parts := strings.Split(strings.Replace(s, "\"", "", -1), ",")
					if absoluteHost {
						var endpoint string
						if host != "" {
							endpoint = strings.TrimRight(host, "/") + ":" + strconv.Itoa(port) + parts[2]
						} else {
							endpoint = parts[2]
						}

						fmt.Fprintln(w, strings.ToUpper(parts[1])+"\t"+endpoint)
					} else {
						fmt.Fprintln(w, strings.Join(parts, "\t"))
					}

				}

				w.Flush()

				fmt.Printf("[%s] Static Assets: \n\n", dirPath)
				// static asset table
				assetMatches := assets.FindAllStringSubmatch(string(dat), -1)
				for _, match := range assetMatches {
					s := strings.TrimPrefix(match[0], "#*")
					parts := strings.Split(s, " ")

					if absoluteHost {
						var endpoint string
						if host != "" {
							endpoint = strings.TrimRight(host, "/") + ":" + strconv.Itoa(port) + strings.TrimLeft(parts[2], ".")
						} else {
							endpoint = parts[2]
						}

						fmt.Fprintln(w, strings.ToUpper(parts[1])+"\t"+endpoint)
					} else {
						fmt.Fprintln(w, strings.Join(parts, "\t"))
					}
				}
				w.Flush()

				// @TODO: need to deal with mounting and static file routers

			}
		}
	}

}
