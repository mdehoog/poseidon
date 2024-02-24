package generator

import (
	"io"
	"text/template"

	"github.com/mdehoog/poseidon/constants"
)

func generateTemplate(w io.Writer, constants *constants.Strings) error {
	t, err := template.New("t").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, constants)
}

const (
	tmpl = `
    package constants

	func init() {
		register(&Strings{
			F: "{{ .F }}",
			C: [][]string{ {{ range $element := .C }}
				{ {{ range $subelement := $element}}
					"{{ $subelement }}",{{ end }}
				},{{ end }}
			},
			M: [][][]string{ {{ range $element := .M }}
				{ {{ range $subelement := $element}}
					{ {{ range $subsub := $subelement}}
						"{{ $subsub }}",{{end}}
					},{{ end }}
				},{{ end }}
			},
			P: [][][]string{ {{ range $element := .P }}
				{ {{ range $subelement := $element}}
					{ {{ range $subsub := $subelement}}
						"{{ $subsub }}",{{end}}
					},{{ end }}
				},{{ end }}
			},
			S: [][]string{ {{ range $element := .S }}
				{ {{ range $subelement := $element}}
					"{{ $subelement }}",{{ end }}
				},{{ end }}
			},
		})
	}
    `
)
