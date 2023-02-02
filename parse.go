package badgrlib

import (
	"strings"

	"github.com/BurntSushi/toml"
)

type PaperFit struct {
	X int `toml:"x"`
	Y int `toml:"y"`
}

type Dimensions struct {
	Width  float64 `toml:"width"`
	Height float64 `toml:"height"`
}

type Object struct {
	FieldName string  `toml:"name"`
	X         float64 `toml:"x"`
	Y         float64 `toml:"y"`
	FontSize  float64 `toml:"font_size"`
}

type Format struct {
	Dimensions Dimensions `toml:"dimensions"`
	PaperFit   PaperFit   `toml:"paper_fit"`
	Objects    []Object   `toml:"objects"`
}

func ParseFormat(src string) (Format, error) {
	var format Format
	_, err := toml.Decode(src, &format)

	if err != nil {
		return Format{}, err
	}

	return format, nil
}

type InputTable struct {
	Data []map[string]string
}

func ParseTable(src string) InputTable {
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")

	fields := strings.Split(lines[0], ";")

	data := make([]map[string]string, 0, len(lines)-1)

	for i := 1; i < len(lines); i++ {
		loop_data := make(map[string]string)

		linesplit := strings.Split(lines[i], ";")

		if len(fields) != len(linesplit) {
			continue
		}

		for idx, field := range fields {
			loop_data[field] = linesplit[idx]
		}

		data = append(data, loop_data)
	}

	return InputTable{Data: data}
}
