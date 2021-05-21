package services

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/PuerkitoBio/goquery"
)

type subSceneItem struct {
	id              string
	title           string
	language        string
	hearingImpaired bool
	owner           string
	comment         string
}

type SubScene struct {
	client http.Client
}

func NewSubSceneService() *SubScene {
	return &SubScene{
		client: http.Client{},
	}
}

func (s *SubScene) Search(kw string) []SubEntry {
	requestBody := strings.NewReader("query=" + kw)
	response, err := s.client.Post("https://subscene.com/subtitles/searchbytitle", "application/x-www-form-urlencoded", requestBody)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode > 300 {
		panic(fmt.Sprintf("Unexpected status code: %d", response.StatusCode))
	}
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	result := []SubEntry{}
	document.Find(".search-result .title a").Each(func(i int, s *goquery.Selection) {
		fmt.Println(">", s.Text())
		result = append(result, SubEntry{
			DisplayName: s.Text(),
			Id:          s.AttrOr("href", ""),
			internal:    nil,
		})
	})

	return result
}

func (s *SubScene) Download(id string) {
	response, err := s.client.Get("https://subscene.com" + id)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	document, _ := goquery.NewDocumentFromReader(response.Body)
	entries := []subSceneItem{}
	document.Find(".content tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.Children().Length() == 1 {
			return
		}

		item := subSceneItem{
			id:              s.Find("a").AttrOr("href", ""),
			language:        strings.TrimSpace(s.Find("td.a1 span").Eq(0).Text()),
			title:           strings.TrimSpace(s.Find("td.a1 span").Eq(1).Text()),
			hearingImpaired: s.Find(".a41").Length() > 0,
			owner:           strings.TrimSpace(s.Find(".a5").Text()),
			comment:         strings.TrimSpace(s.Find(".a6").Text()),
		}
		entries = append(entries, item)
	})
	languages := []string{}
	{
		langMap := map[string]bool{}
		for _, v := range entries {
			if _, found := langMap[v.language]; !found {
				languages = append(languages, v.language)
				langMap[v.language] = true
			}
		}
	}
	options := languages
	answer := 0
	survey.AskOne(&survey.Select{
		Message: "Choose your prefered language",
		Options: options,
	}, &answer, nil)
	entries = pickByLang(entries, options[answer])
	options = []string{}
	for _, entry := range entries {
		options = append(options, entry.title)
	}
	survey.AskOne(&survey.Select{
		Message: "Pick a title",
		Options: options,
	}, &answer, nil)
	s.DownloadStage2(entries[answer].id)
}

func pickByLang(entries []subSceneItem, lang string) []subSceneItem {
	result := []subSceneItem{}
	for _, item := range entries {
		if item.language == lang {
			result = append(result, item)
		}
	}

	return result
}

func (s *SubScene) DownloadStage2(id string) {
	response, err := s.client.Get("https://subscene.com" + id)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	document, _ := goquery.NewDocumentFromReader(response.Body)
	href, _ := document.Find(".download a").Attr("href")
	href = "https://subscene.com" + href
	zipFile, err := os.CreateTemp(os.TempDir(), "subscene-*.zip")
	if err != nil {
		panic(err)
	}
	response, err = s.client.Get(href)
	if err != nil {
		panic(err)
	}
	io.Copy(zipFile, response.Body)
	zipFile.Close()
	zipFile, _ = os.Open(zipFile.Name())
	archive, _ := zip.OpenReader(zipFile.Name())
	for _, f := range archive.File {
		fmt.Println("Found file", f.Name)
		dst, err := os.Create("/tmp/bach.srt")
		if err != nil {
			panic(err)
		}
		src, err := f.Open()
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			panic(err)
		}
		dst.Close()
		src.Close()
		fmt.Println("File written")
	}
}
