package badgrlib

import (
	"strings"

	"github.com/BurntSushi/toml"
)

type Dimensions struct {
	Width  float32 `toml:"width"`
	Height float32 `toml:"height"`
}

type Object struct {
	FieldName string  `toml:"name"`
	X         float32 `toml:"x"`
	Y         float32 `toml:"y"`
	Width     float32 `toml:"w"`
	Height    float32 `toml:"h"`
}

type Format struct {
	Dimensions Dimensions `toml:"dimensions"`
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

	data := make([]map[string]string, len(lines)-1)

	for i := 1; i < len(lines); i++ {
		data[i-1] = make(map[string]string)

		linesplit := strings.Split(lines[i], ";")

		if len(fields) != len(linesplit) {
			continue
		}

		for idx, field := range fields {
			data[i-1][field] = linesplit[idx]
		}
	}

	return InputTable{Data: data}
}
