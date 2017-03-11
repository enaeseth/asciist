package service

import (
	"bytes"
	"fmt"
	"image"
	"net/http"

	// Image decoders
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/enaeseth/asciist/convert"
	"gopkg.in/gin-gonic/gin.v1"
)

func SetDebug(debug bool) {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func New() http.Handler {
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		var req Request

		if err := c.BindJSON(&req); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}

		img, err := decodeImage(&req)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}

		art := convert.FromImage(img, req.Width)
		c.JSON(http.StatusOK, Success{
			Art: art.String(),
		})
	})

	router.NoRoute(func(c *gin.Context) {
		fail(c, http.StatusNotFound, fmt.Errorf("not found: %s", c.Request.URL.Path))
	})

	return router
}

func fail(c *gin.Context, code int, err error) {
	c.JSON(code, Failure{
		Error: err.Error(),
	})
}

func decodeImage(req *Request) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(req.Image))
	return img, err
}
