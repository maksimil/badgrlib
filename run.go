package badgrlib

import (
	"sync"
)

func createSvgs(format Format, data_source string) ([]string, error) {
	data := ParseTable(data_source)

	var mut sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(data.Data))

	svgs := make([]string, 0, len(data.Data))
	ret_err := error(nil)

	for _, loop_object_data := range data.Data {
		go func(object_data map[string]string) {
			svg, err := CreateSingleSvg(format, object_data)

			mut.Lock()
			if err != nil {
				ret_err = err
			} else {
				svgs = append(svgs, svg)
			}
			mut.Unlock()

			wg.Done()
		}(loop_object_data)
	}

	wg.Wait()

	if ret_err != nil {
		return []string{}, ret_err
	}

	return svgs, nil
}

func createPapers(format Format, svgs []string) ([]string, error) {
	objects_on_page := format.PaperFit.X * format.PaperFit.Y

	var mut sync.Mutex

	var wg sync.WaitGroup

	papers := make([]string, 0)
	ret_err := error(nil)

	loop_svg_group := make([]string, 0, objects_on_page)

	create_paper := func() {
		wg.Add(1)
		go func(svg_group []string) {
			paper, err := FitSvgsToPaper(format, svg_group)

			mut.Lock()
			if err != nil {
				ret_err = err
			} else {
				papers = append(papers, paper)
			}
			mut.Unlock()

			wg.Done()
		}(loop_svg_group)
	}

	for _, svg := range svgs {
		loop_svg_group = append(loop_svg_group, svg)

		if len(loop_svg_group) == objects_on_page {
			create_paper()
			loop_svg_group = make([]string, 0, objects_on_page)
		}
	}

	if len(loop_svg_group) > 0 {
		create_paper()
	}

	wg.Wait()

	if ret_err != nil {
		return []string{}, ret_err
	}

	return papers, nil
}

func RunFormat(format_source string, data_source string) ([]string, error) {
	format, err := ParseFormat(format_source)

	if err != nil {
		return []string{}, err
	}

	svgs, err := createSvgs(format, data_source)

	if err != nil {
		return []string{}, err
	}

	papers, err := createPapers(format, svgs)

	if err != nil {
		return []string{}, err
	}

	return papers, nil
}
