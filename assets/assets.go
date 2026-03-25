package assets

import (
    "embed"
    "io/fs"
)

//go:embed templates
var all embed.FS

var TemplatesFS fs.FS

func init() {
    sub, err := fs.Sub(all, "templates")
    if err != nil {
        panic(err)
    }
    TemplatesFS = sub
}