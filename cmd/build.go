package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, _ := os.ReadFile("./data/books.json")
		var books []Book
		json.Unmarshal(data, &books)

		for i, b := range books {
			slog.Info("Fetching metadata for", "title", b.Title)
			img, desc, _ := FetchAmazonMetadata(b.AmazonLink)
			books[i].Image = img
			books[i].Description = desc
		}

		litter.Dump(books)

		tmpl, err := os.ReadFile("old.html")
		if err != nil {
			panic(err)
		}

		funcs := template.FuncMap{
			"StarRating": func(r int) string {
				return strings.Repeat("★", r) + strings.Repeat("☆", 5-r)
			},
		}

		t := template.Must(template.New("page").Funcs(funcs).Parse(string(tmpl)))
		f, _ := os.Create("index.html")
		defer f.Close()
		t.Execute(f, map[string]interface{}{
			"Books": books,
		})
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Review      string `json:"review"`
	Rating      int    `json:"rating"`
	AmazonLink  string `json:"amazon_link"`
	Image       string
	Description string
}

func FetchAmazonMetadata(url string) (image, desc string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Bot/1.0)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	image, _ = doc.Find(`meta[property="og:image"]`).Attr("content")
	desc, _ = doc.Find(`meta[property="og:description"]`).Attr("content")
	return image, desc, nil
}
