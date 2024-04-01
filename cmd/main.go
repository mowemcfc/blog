package main

import (
	"html/template"
	"io"
  "log"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
  ID string
  Title string
  Blurb string
  Body string
}

type HomePage struct {
  Posts map[string]Post
}

func newPost(id string, title string, blurb string, body string) Post {
  createdAt := time.Now().UTC().Format(time.RFC1123)
  return Post{
    CreatedAt: createdAt,
    ID: id,
    Title: title,
    Blurb: blurb,
    Body: body,
  }
}

func newHomePage(posts map[string]Post) HomePage {
  return HomePage {
    Posts: posts,
  }
}

type BlogPage struct {
  Post Post
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

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
	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")
  e.Static("/js", "js")

	e.GET("/", func(c echo.Context) error {
    log.Println(homePage.Posts["first"].ID)
		return c.Render(200, "index", homePage)
	})

  e.GET("/blog/:id", func(c echo.Context) error {
    id := c.Param("id")
    post := allBlogPosts[id]
    return c.Render(200, "blog", post)
  })

	// e.POST("/contacts", func(c echo.Context) error {
	// 	name := c.FormValue("name")
	// 	email := c.FormValue("email")
	//
	// 	if page.Data.hasEmail(email) {
	// 		formData := newFormData()
	// 		formData.Values["name"] = name
	// 		formData.Values["email"] = email
	// 		formData.Errors["email"] = "email already exists"
	// 		return c.Render(422, "createContactForm", formData)
	// 	}
	//
	// 	contact := newContact(name, email)
	// 	page.Data.Contacts = append(page.Data.Contacts, contact)
	// 	c.Render(200, "createContactForm", newFormData())
	// 	return c.Render(200, "oob-contact", contact)
	// })
	//
	// e.DELETE("/contacts/:id", func(c echo.Context) error {
	// 	time.Sleep(1 * time.Second)
	// 	idStr := c.Param("id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		return c.String(400, "invalid id")
	// 	}
	// 	index := page.Data.indexOf(id)
	// 	if index == -1 {
	// 		return c.String(404, "contact not found")
	// 	}
	// 	page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)
	//
	// 	return c.NoContent(200)
	// })

	e.Logger.Fatal(e.Start(":42069"))
}
