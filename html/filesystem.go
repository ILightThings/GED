package html

import (
	"embed"
	_ "embed"
)

////go:embed credtable.html
//var CredTable string
//
////go:embed footer.html
//var Footer string
//
////go:embed header.html
//var Header string
//
////go:embed hosttable.html
//var HostTable string
//
////go:embed import.html
//var Import string
//
////go:embed setting.html
//var Setting string
//
////go:embed updateCred.html
//var UpdateCred string
//
////go:embed updateHost.html
//var UpdateHost string

//go:embed *.html
var HTML embed.FS
