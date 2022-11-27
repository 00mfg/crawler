package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	url := "https://www.thepaper.cn/"

	body, err := Fetch(url)

	if err != nil {
		fmt.Printf("read content failed %v", err)
		return
	}

	//res := REParse(body)
	//res := XpathParse(body)
	res := CSSParse(body)
	fmt.Println(res)

}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %v", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DerterminEncoding(bodyReader) //utf-8
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func DerterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	if err != nil {
		fmt.Printf("fetch error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

func REParse(page []byte) []string {
	var res []string
	matches := headerRe.FindAllSubmatch(page, -1)

	for _, m := range matches {
		res = append(res, string(m[1]))
	}
	return res
}

func XpathParse(body []byte) []string {
	var res []string
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Printf(" htmlquery Parse failed: %v", err)
	}

	nodes := htmlquery.Find(doc, `//div[@class="small_toplink__GmZhY"]/a/h2`)

	for _, node := range nodes {
		res = append(res, node.FirstChild.Data)
	}
	return res
}

func CSSParse(body []byte) []string {
	var res []string
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("read content failed: %v", err)
	}
	doc.Find("div a h2").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		res = append(res, title)
	})
	return res
}
