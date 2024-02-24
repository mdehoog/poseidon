package main

import (
	"io"
	"text/template"
)

const (
	tmpl = `
    package poseidon

    type constantsStr struct {
        C [][]string
        S [][]string
        M [][][]string
        P [][][]string
    }

	//nolint:lll,dupl // long lines, duplicated parts
    var cs = constantsStr {
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
    }
    `
)

type constants struct {
	C [][]string   // compRoundConstants
	S [][]string   // sparseMatrix
	M [][][]string // MDS matrix
	P [][][]string // preSparseMatrix
}

func GenerateTemplate(w io.Writer, constants *constants) error {
	t, err := template.New("t").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, constants)
}
