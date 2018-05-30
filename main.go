package main

import (
  "encoding/csv"
  "fmt"
  "net/url"
  "os"
  "strconv"
  "strings"

  "github.com/jcelliott/lumber"
  "gopkg.in/urfave/cli.v1" // imports as package "cli"
)

type redirectRequest struct {
  index int
  country string
  locale string
  host string
  path string
  target string
}

func main() {
  lumber.Debug("starting")

  app := cli.NewApp()
  app.Name = "csv2geo"
  app.Usage = "generate georedirect rules from CSV"
  app.Version = "0.1.0"

  app.Action = func(c *cli.Context) error {
    arg := c.Args().Get(0)

    err := doAction(arg)

    return err
  }


  err := app.Run(os.Args)
  if err != nil {
    os.Exit(1)
  }

  return
}

func doAction(arg string) error {
  // check whether file exists
  file, err := os.Open(arg);
  if err != nil {
    return err
  }
  defer file.Close()

  // parse CSV from file
  r := csv.NewReader(file)
  lines, err := r.ReadAll()

  // process each line
  for _, line := range lines {
    parsed, err := parseLine(line)
    if err == nil {
      formatLine(parsed)
    }
  }

  return err
}

func parseLine(line []string) (redirectRequest, error) {
  index, err := strconv.Atoi(line[0])

  if err != nil {
    return redirectRequest{}, err
  }

  url, err := url.Parse(fmt.Sprintf("http://%s", line[6]))

  if err != nil {
    return redirectRequest{}, err
  }

  parsed := redirectRequest{index: index, country: line[3], locale: line[4], host: url.Hostname(), path: url.EscapedPath(), target: line[8]}

  return parsed, err
}

func formatLine(r redirectRequest) {
  geocode := strings.ToUpper(r.country)
  path := r.path
  target := r.target

  fmt.Printf("GEO:%s^%s(.*) %s$1 [R,L]\n", geocode, path, target)
}
