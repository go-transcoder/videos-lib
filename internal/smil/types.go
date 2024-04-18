package smil

import "encoding/xml"

// SMIL structure
type SMIL struct {
	XMLName xml.Name `xml:"smil"`
	Title   string   `xml:"title,attr"`
	Body    Body     `xml:"body"`
}

// Body structure
type Body struct {
	Switch Switch `xml:"switch"`
}

// Switch structure
type Switch struct {
	Videos []Video `xml:"video"`
}

// Video structure
type Video struct {
	Src            string `xml:"src,attr"`
	Height         string `xml:"height,attr"`
	Width          string `xml:"width,attr"`
	SystemLanguage string `xml:"systemLanguage,attr"`
	SystemBitrate  string `xml:"system-bitrate,attr"`
}
