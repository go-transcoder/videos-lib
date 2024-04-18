package smil

import (
	"encoding/xml"
	"os"
)

type Command func(string) error

func (Command Command) Exec(outputPath string) error {
	// creating the smil file content
	smil := SMIL{
		Body: Body{
			Switch: Switch{
				Videos: []Video{
					{
						Src:            "240p/video.mp4",
						Height:         "240",
						Width:          "424",
						SystemLanguage: "eng",
						SystemBitrate:  "450000",
					},
				},
			},
		},
	}

	// Encode SMIL structure to XML
	xmlData, err := xml.MarshalIndent(smil, "", "    ")
	if err != nil {
		panic(err)
	}

	// Write the XML data to a file
	file, err := os.Create(outputPath + "/video.smil")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(xmlData)

	return nil
}
