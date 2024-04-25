package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Post struct {
	CreatedAt string
	ID        string
	Title     string
	Blurb     string
	Body      string
}

type HomePage struct {
	Posts map[string]Post
}

func newPost(id string, title string, blurb string, body string) Post {
	createdAt := time.Now().UTC().Format(time.RFC1123)
	return Post{
		CreatedAt: createdAt,
		ID:        id,
		Title:     title,
		Blurb:     blurb,
		Body:      body,
	}
}

func newHomePage(posts map[string]Post) HomePage {
	return HomePage{
		Posts: posts,
	}
}

type BlogPage struct {
	Post Post
}

var adapter *echoadapter.EchoLambda

func init() {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Renderer = newTemplate()

	allBlogPosts := map[string]Post{
		"first": newPost(
			"first",
			"First Blog",
			"My first blog...",
			"Just testing my first blog post here",
		),
		"second": newPost(
			"second",
			"Second Blog",
			"My second blog...",
			"Just testing my second blog post here",
		),
		"third": newPost(
			"third",
			"Third Blog",
			"My third blog...",
			"Just testing my third blog post here",
		),
	}

	homePage := newHomePage(allBlogPosts)
	app.Static("/images", "images")
	app.Static("/css", "css")
	app.Static("/js", "js")

	app.GET("/", func(c echo.Context) error {
		err := c.Render(http.StatusOK, "index", homePage)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})

	app.GET("/blog/:id", func(c echo.Context) error {
		id := c.Param("id")
		post := allBlogPosts[id]
		err := c.Render(http.StatusOK, "blog", post)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})

	isLambda := os.Getenv("LAMBDA_TASK_ROOT") != ""
	if isLambda {
		adapter = echoadapter.New(app)
	} else {
		log.Fatal(app.Start(":42069"))
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
