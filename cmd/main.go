package main

import (
  "html/template"
  "io"

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

type Count struct {
  Count int
}

type Contact struct {
  Name string
  Email string
}

type Contacts = []Contact

func newContact(name string, email string) Contact {
  return Contact{
    Name: name,
    Email:email,
  }
}

type Data struct {
  Contacts Contacts
}
func newData() Data {
  return Data{
    Contacts: []Contact{
      newContact("John", "jd@gmail.com"),
      newContact("Glenda", "gd@gmail.com"),
    },
  }
}


func main() {
  e := echo.New()
  e.Use(middleware.Logger())

  data := newData()
  e.Renderer = newTemplate() 

  e.GET("/", func(c echo.Context) error {
    return c.Render(200, "index", data)
  })

  e.POST("/contacts", func(c echo.Context) error {
    name := c.FormValue("name")
    email := c.FormValue("email")

    data.Contacts = append(data.Contacts, newContact(name, email))
    return c.Render(200, "contacts", data)
  })

  e.Logger.Fatal(e.Start(":42069"))
}
