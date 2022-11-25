package html

import (
	"github.com/mrasu/GravityR/injection"
	"github.com/pkg/errors"
	"html/template"
	"io/fs"
	"os"
)

const (
	TypeMain    = "main"
	TypeMermaid = "mermaid"
)

func CreateHtml(htmlType, filename string, bo *BuildOption) error {
	script, err := fs.ReadFile(injection.ClientDist, "client/dist/assets/"+htmlType+".js")
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	style, err := fs.ReadFile(injection.ClientDist, "client/dist/assets/"+htmlType+".css")
	if err != nil {
		return errors.Wrap(err, "failed to open file")
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
		return errors.Wrap(err, "failed create template")
	}
	m := map[string]interface{}{
		"script": template.JS(script),
		"style":  template.CSS(style),
		"gr":     bo.createGrMap(),
	}

	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create html file")
	}

	err = tpl.Execute(f, m)
	if err != nil {
		return errors.Wrap(err, "failed to write template to html")
	}

	return nil
}
