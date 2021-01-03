package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	escli "github.com/KanchiShimono/elasticsearch-client-examples/client"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/urfave/cli/v2"
)

func search(q string, s int64) error {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}
	client := escli.NewClient(es)

	info, err := client.ES.Info()
	if err != nil {
		return err
	}
	defer info.Body.Close()

	log.Println(info)

	results, err := client.Search(q, s)
	if err != nil {
		return err
	}
	for _, r := range results {
		log.Printf("Title=%s, Score=%f", r.Title, r.Score)
	}

	return nil
}

func main() {
	c := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "search",
				Description: "search",
				ArgsUsage:   "query size",
				Action: func(c *cli.Context) error {
					q := c.Args().Get(0)
					s, err := strconv.ParseInt(c.Args().Get(1), 10, 64)
					if err != nil {
						return errors.New("size is required as int")
					}
					return search(q, s)
				},
			},
		},
	}

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
