// Copyright 2020 Kentaro Hibino. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/cloudjjcc/asynq/internal/rdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Shows all running worker servers",
	Long: `Servers (asynq servers) will show all running worker servers
pulling tasks from the specified redis instance.

The command shows the following for each server:
* Host and PID of the process in which the server is running
* Number of active workers out of worker pool
* Queue configuration
* State of the worker server ("running" | "quiet")
* Time the server was started

A "running" server is pulling tasks from queues and processing them.
A "quiet" server is no longer pulling new tasks from queues`,
	Args: cobra.NoArgs,
	Run:  servers,
}

func init() {
	rootCmd.AddCommand(serversCmd)
}

func servers(cmd *cobra.Command, args []string) {
	r := rdb.NewRDB(redis.NewClient(&redis.Options{
		Addr:     viper.GetString("uri"),
		DB:       viper.GetInt("db"),
		Password: viper.GetString("password"),
	}))

	servers, err := r.ListServers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(servers) == 0 {
		fmt.Println("No running servers")
		return
	}

	// sort by hostname and pid
	sort.Slice(servers, func(i, j int) bool {
		x, y := servers[i], servers[j]
		if x.Host != y.Host {
			return x.Host < y.Host
		}
		return x.PID < y.PID
	})

	// print server info
	cols := []string{"Host", "PID", "State", "Active Workers", "Queues", "Started"}
	printRows := func(w io.Writer, tmpl string) {
		for _, info := range servers {
			fmt.Fprintf(w, tmpl,
				info.Host, info.PID, info.Status,
				fmt.Sprintf("%d/%d", info.ActiveWorkerCount, info.Concurrency),
				formatQueues(info.Queues), timeAgo(info.Started))
		}
	}
	printTable(cols, printRows)
}

// timeAgo takes a time and returns a string of the format "<duration> ago".
func timeAgo(since time.Time) string {
	d := time.Since(since).Round(time.Second)
	return fmt.Sprintf("%v ago", d)
}

func formatQueues(qmap map[string]int) string {
	// sort queues by priority and name
	type queue struct {
		name     string
		priority int
	}
	var queues []*queue
	for qname, p := range qmap {
		queues = append(queues, &queue{qname, p})
	}
	sort.Slice(queues, func(i, j int) bool {
		x, y := queues[i], queues[j]
		if x.priority != y.priority {
			return x.priority > y.priority
		}
		return x.name < y.name
	})

	var b strings.Builder
	l := len(queues)
	for _, q := range queues {
		fmt.Fprintf(&b, "%s:%d", q.name, q.priority)
		l--
		if l > 0 {
			b.WriteString(" ")
		}
	}
	return b.String()
}
