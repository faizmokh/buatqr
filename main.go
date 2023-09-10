// QR code generator

// Requirements
// 1. Display a page with a form to enter a URL
// 2. Display a QR code for the URL entered in the form
// 3. Allow the user to save the QR code as a PNG file

package main

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.POST("/generate", func(c *gin.Context) {
		url := c.PostForm("link")

		if !isValidURL(url) {
			c.HTML(200, "index.html", gin.H{
				"error": "Boleh letak url je",
			})
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/generate?link="+url)
	})

	router.GET("/generate", func(c *gin.Context) {
		url := c.Query("link")
		png, err := qrcode.Encode(url, qrcode.Medium, 512)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		encodedPng := base64.StdEncoding.EncodeToString(png)
		qrImgSafeHTML := template.HTML(encodedPng)

		c.HTML(200, "index.html", gin.H{
			"url":    url,
			"qrcode": qrImgSafeHTML,
		})
	})

	router.Run("localhost:8080")
}

func isValidURL(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
