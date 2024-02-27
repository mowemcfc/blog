package main

import (
	"html/template"
	"io"
  "time"
	// "strconv"
	// "time"

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
  CreatedAt time.Time
  Title string
  Blurb string
  Body string
}
type Posts = []Post

type HomePage struct {
  Posts Posts
}

func newPost(title string, blurb string, body string) Post {
  // TODO: decide on localisation approach
  createdAt := time.Now().UTC()
  return Post{
    CreatedAt: createdAt,
    Title: title,
    Blurb: blurb,
    Body: body,
  }
}

func newHomePage() HomePage {
  return HomePage {
    Posts: Posts{
      newPost("title1", "blurb1", "body1"),
      newPost("title2", "blurb2", "body2"),
    },
  }
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	homePage := newHomePage()
	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", homePage)
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
