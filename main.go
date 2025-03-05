package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"urlshortener/db"
	"urlshortener/frontend"

	"github.com/PuerkitoBio/purell"
	_ "github.com/joho/godotenv/autoload"
)

var ctx = context.Background()
var config *Config
var executor *db.Queries

//go:embed assets/*
var assets embed.FS

func main() {
	config = GetConfig()

	connector, dir := ConnectToDatabase(config)
	if dir != "" {
		fmt.Println("Connected to turso")
		defer os.RemoveAll(dir)
	} else {
		fmt.Println("Connected to local.db")
	}

	executor = db.New(connector)
	defer connector.Close()

	router := http.NewServeMux()

	router.HandleFunc("GET /", GzipF(home))
	router.HandleFunc("GET /custom_code", func(w http.ResponseWriter, r *http.Request) {
		frontend.CustomCode().Render(ctx, w)
	})
	router.HandleFunc("DELETE /custom_code", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("GET /{code}", shortLink)
	router.HandleFunc("POST /new_link", newLink)
	router.Handle("GET /assets/", Gzip(http.FileServer(http.FS(assets))))

	server := http.Server{
		Addr:    config.ListenAddress,
		Handler: router,
	}
	fmt.Println("Starting server at", "http://"+config.FQDN)
	log.Fatal(server.ListenAndServe())
}
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	frontend.HomePage().Render(ctx, w)
}
func shortLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.NotFound(w, r)
		return
	}
	long, err := executor.GetFromCode(ctx, code)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	fmt.Printf("%s accessed %s: X-Forwarded-For: %s\n", r.RemoteAddr, code, r.Header.Get("X-Forwarded-For"))
	http.Redirect(w, r, long, http.StatusTemporaryRedirect)
}
func newLink(w http.ResponseWriter, r *http.Request) {
	link := r.FormValue("link")
	customCode := strings.TrimSpace(r.FormValue("code"))

	formattedLink := link
	if !strings.HasPrefix(link, "https://") && !strings.HasPrefix(formattedLink, "http://") {
		formattedLink = "https://" + link
	}
	formattedLink, err := purell.NormalizeURLString(formattedLink, purell.FlagLowercaseScheme|purell.FlagLowercaseHost|purell.FlagUppercaseEscapes|purell.FlagRemoveTrailingSlash)
	if err != nil {
		fmt.Fprint(w, "Invalid link, failed to format page")
		return
	}

	if customCode != "" {
		if len(customCode) > 16 || len(customCode) < 3 {
			fmt.Fprint(w, "Custom code must be between 3 and 16")
			return
		}
		existingLink, err := executor.GetFromCode(ctx, customCode)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err.Error())
			fmt.Fprint(w, "Something went wrong")
			return
		}
		if existingLink == formattedLink {
			newLink := generateLink(config, customCode)
			frontend.Link(newLink).Render(ctx, w)
			return
		}
		if existingLink != "" {
			fmt.Fprint(w, "Code is being used")
			return
		}
	}

	code, _ := executor.GetFromUrl(ctx, formattedLink)
	if code != "" {
		newLink := generateLink(config, code)
		frontend.Link(newLink).Render(ctx, w)
		return
	}
	res, err := http.Get(formattedLink)
	if err != nil {
		fmt.Fprint(w, "Invalid link, failed to get page")
		return
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		fmt.Fprint(w, "Invalid link, incorrect status ", res.StatusCode)
		return
	}
	if customCode != "" {
		code = customCode
	} else {
		code = createCode(5)
		for code, _ := executor.GetFromCode(ctx, code); code != ""; {
			code = createCode(5)
		}
	}
	executor.CreateLink(ctx, db.CreateLinkParams{Code: code, LongUrl: formattedLink})
	newLink := generateLink(config, code)
	frontend.Link(newLink).Render(ctx, w)
}

var charset = strings.Split("abcdefghijklmnopqrstuvwxyz0123456789", "")

func createCode(length int) string {
	var code string
	for range length {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatal(err.Error())
		}
		code += charset[idx.Int64()]
	}
	return code
}
func generateLink(config *Config, code string) string {
	return "http://" + config.FQDN + "/" + code
}
