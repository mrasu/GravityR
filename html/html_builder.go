package html

import (
	"github.com/mrasu/GravityR/injection"
	"html/template"
	"io/fs"
	"os"
)

func CreateHtml(filename string, bo *BuildOption) error {
	script, err := fs.ReadFile(injection.ClientDist, "client/dist/assets/main.js")
	if err != nil {
		return err
	}
	style, err := fs.ReadFile(injection.ClientDist, "client/dist/assets/main.css")
	if err != nil {
		return err
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>GravityR</title>
	<script>
      window.grParam = {{.gr}}
    </script>
	<style>{{.style}}</style>
  </head>
  <body>
    <div id="app"></div>
    <script>{{.script}}</script>
  </body>
</html>
`
	tpl, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	m := map[string]interface{}{
		"script": template.JS(script),
		"style":  template.CSS(style),
		"gr":     bo.createGrMap(),
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = tpl.Execute(f, m)
	if err != nil {
		return err
	}

	return nil
}
