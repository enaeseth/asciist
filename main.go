package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/enaeseth/asciist/client"
	"github.com/enaeseth/asciist/convert"
	"github.com/enaeseth/asciist/service"
)

var (
	app = kingpin.New("asciist", "ASCII art as a service")

	cli         = app.Command("convert", "Convert an image using the service")
	clientURL   = cli.Flag("url", "URL of the service").Short('u').Default("http://localhost:27244").String()
	clientWidth = cli.Flag("width", "Character width to output").Short('w').Uint()
	clientFile  = cli.Arg("file", "Input image file").File()

	serve = app.Command("serve", "Run the service")
	debug = serve.Flag("debug", "Use debugging mode").Bool()
	host  = serve.Flag("host", "Interface to listen on").Envar("HOST").String()
	port  = serve.Flag("port", "Port to listen on").Default("27244").Envar("PORT").Int()

	test      = app.Command("test", "Convert an image locally")
	testWidth = test.Flag("width", "Character width to output").Short('w').Uint()
	testFile  = test.Arg("file", "Input image file").File()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cli.FullCommand():
		runClient()

	case serve.FullCommand():
		runService()

	case test.FullCommand():
		runLocal()
	}
}

func runClient() {
	c := client.New(*clientURL)
	art, err := c.Convert(readInput(clientFile), width(clientWidth))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(art)
}

func runService() {
	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	listen := fmt.Sprintf("%s:%d", *host, *port)
	s := service.New()

	fmt.Printf("Listening on %s\n", listen)
	s.Run(listen)
}

func runLocal() {
	img, _, err := image.Decode(readInput(testFile))
	if err != nil {
		log.Fatal(err)
	}

	art := convert.FromImage(img, width(testWidth))
	fmt.Println(art)
}

func readInput(fileArg **os.File) *os.File {
	if *fileArg != nil {
		return *fileArg
	}

	return os.Stdin
}

func width(widthArg *uint) uint {
	if widthArg != nil && *widthArg != 0 {
		return *widthArg
	}

	return defaultWidth()
}
