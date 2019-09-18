// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/wfscheper/xkcdpwd/internal/langs"
)

func main() {
	err := vfsgen.Generate(langs.Languages, vfsgen.Options{
		PackageName:  "langs",
		BuildTags:    "!dev",
		VariableName: "Languages",
	})
	if err != nil {
		log.Fatalln(err)
	}

}
