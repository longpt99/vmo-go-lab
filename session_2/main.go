package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"regexp"

	"golang.org/x/net/html"
)

func main() {
	// lesson1()
	// lesson2()
	// lesson3()
}

func lesson1() {
	input := "This is a sample text. This text is just an example."
	response := make(map[string]int)
	arrInput := strings.Split(input, " ")

	for _, text := range arrInput {
		textFormatted := strings.ReplaceAll(text, ".", "")
		val := response[textFormatted]
		response[textFormatted] = val + 1
	}

	for key, res := range response {
		fmt.Printf("%s: %d\n", key, res)
	}
}

func lesson2() {
	var num int
	var arrNum []int = []int{2, 3, 5, 4, 7, 8, 11, 21, 17, 41, 40, 42, 43}

	for {
		fmt.Print("Nhap so nguyen nguyen duong: ")
		fmt.Scan(&num)

		if num >= 0 {
			break
		}
	}

	for _, n := range arrNum {
		isValid := isPrimeNumber((n))

		if isValid && n <= num {
			fmt.Printf("%d ", n)
		}
	}

}

func isPrimeNumber(n int) bool {
	if n < 2 {
		return false
	} else {
		for i := 2; i <= n-1; i++ {
			if n%i == 0 {
				return false
			}
		}
	}

	return true
}

func lesson3() {
	resp, err := http.Get("https://tinhte.com")

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		log.Panic(err)
	}

	arrText := getText(doc)

	re := regexp.MustCompile(`(?i)https?://[\w\-\.]+\.\w+(:\d+)?(/[\w\-\./\?%&=]*)?`)

	var urls []string
	for _, t := range arrText {
		/*
			This statement uses a regular expression to find all the URLs that are present in the text content. The regular expression is defined in the re variable and matches URLs of the form http://, https://, or www..
			The -1 argument means that all matches should be returned
		*/
		urls = append(urls, re.FindAllString(t, -1)...)
	}

	fmt.Println(urls)

}

func getText(n *html.Node) []string {
	/*
		This statement checks if the current node is a text node.
		If it is, the function returns the data of the text node as a slice of strings.
	*/
	if n.Type == html.TextNode {
		return []string{n.Data}
	}
	/*
		This statement checks if the current node is a script tag.
		If it is, the function returns nil, which means that the script tag is ignored.
	*/
	if n.Type == html.ElementNode && n.Data == "script" {
		return nil // ignore script tags
	}
	var text []string

	/*
		This statement iterates over the child nodes of the current node
	*/
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text = append(text, getText(c)...)
	}
	return text
}
