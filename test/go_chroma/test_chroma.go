package test

import (
	"bytes"
	"log"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func Go_chroma() {
	// s := fmt.Sprintf("./go_chroma/%s", "test_chroma.go")

	// content, err := os.ReadFile(s)

	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(string(content))
	long_str := `[print("hello world") for i in range(100)]`

	lexer := lexers.Analyse(long_str)

	if lexer == nil {
		lexer = lexers.Fallback
	}

	// fmt.Println(lexer)

	// formatter := html.New(html.WithLineNumbers(true), html.Standalone(true))

	formatter := formatters.Get("html")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, long_str)
	if err != nil {
		log.Println(err)
	}

	var buf bytes.Buffer

	err = formatter.Format(&buf, styles.GitHub, iterator)
	if err != nil {
		log.Println(err)
	}

	os.WriteFile("./tmp/test_chroma_v1.html", buf.Bytes(), os.ModePerm)

}
