package cli

import (
	"fmt"
	"strings"
	"tank/core"

	flag "github.com/spf13/pflag"
)

var topImg, bottomImg, outputImg string

const defaultOutputImg = "output.png"

func init() {
	flag.StringVarP(&topImg, "top", "t", "", "Path to top image")
	flag.StringVarP(&bottomImg, "bottom", "b", "", "Path to bottom image")
	flag.StringVarP(&outputImg, "output", "o", defaultOutputImg, "Path to output image")
	flag.CommandLine.SortFlags = false
}

func Run() {
	flag.Parse()
	if topImg == "" || bottomImg == "" {
		fmt.Println("missing required arguments: --top, --bottom")
		return
	}
	if !strings.HasSuffix(outputImg, ".png") {
		outputImg += ".png"
	}
	if err := tank.Make(topImg, bottomImg, outputImg); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Output image saved to %s\n", outputImg)
}
