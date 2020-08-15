package cmd

import (
	"fmt"
	"html/template"
	"os"
	"testing"

	//"github.com/spf13/cobra"
	"github.com/richarda23/rspace-client-go/rspace"
)

type Weather struct {
	Temp int
	Wind template.HTML
}

func (w *Weather) Temp2() template.HTML {
	return template.HTML(fmt.Sprintf("<fileId=%d>", w.Temp))
}

var templ = ` {{range $val := .}} 
 <br> Test template parse: {{$val.Temp2}}, wind is {{$val.Wind}}  </br>
  {{end}}
  `

func TestBasicTemplate(t *testing.T) {
	weather := Weather{23, template.HTML(`<fileId="windy">`)}
	weather2 := Weather{28, template.HTML(`<fileId="Calm">`)}

	t2 := template.Must(template.New("templ").Parse(templ))
	t2.Execute(os.Stderr, []Weather{weather, weather2})
}
func TestTemplate(t *testing.T) {

	info := &rspace.IdentifiableNamable{}
	info.Id = 1234
	info.GlobalId = "GL1234"
	info.Name = "test"
	fInfo := rspace.FileInfo{IdentifiableNamable: info}
	results := make([]*rspace.FileInfo, 0)
	results = append(results, &fInfo)
	_, err := generateSummaryContent(results)
	if err != nil {
		t.Fatalf(err.Error())
	}
	//	messageStdErr(html)
	// test-spy for context error-writer
}
