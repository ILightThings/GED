package html

import (
	"embed"
	_ "embed"
)

////go:embed table_cred.html
//var CredTable string
//
////go:embed footer.html
//var Footer string
//
////go:embed header.html
//var Header string
//
////go:embed table_host.html
//var HostTable string
//
////go:embed import.html
//var Import string
//
////go:embed setting.html
//var Setting string
//
////go:embed updateform_cred.html
//var UpdateCred string
//
////go:embed updateform_host.html
//var UpdateHost string

//go:embed *.html
var HTML embed.FS
