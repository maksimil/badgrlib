package badgrlib

import (
	"sync"
)

func createSvgs(format Format, data_source string) ([]string, error) {
	data := ParseTable(data_source)

	var mut sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(data.Data))

	svgs := make([]string, len(data.Data))
	ret_err := error(nil)

	for loop_idx, loop_object_data := range data.Data {
		go func(idx int, object_data map[string]string) {
			svg, err := CreateSingleSvg(format, object_data)

			if err != nil {
				mut.Lock()
				ret_err = err
				mut.Unlock()
			} else {
				svgs[idx] = svg
			}

			wg.Done()
		}(loop_idx, loop_object_data)
	}

	wg.Wait()

	if ret_err != nil {
		return []string{}, ret_err
	}

	return svgs, nil
}

func createPages(format Format, svgs []string) ([]string, error) {
	objects_on_page := format.PaperFit.X * format.PaperFit.Y

	page_count := len(svgs) / objects_on_page
	if len(svgs)%objects_on_page > 0 {
		page_count += 1
	}

	var mut sync.Mutex

	var wg sync.WaitGroup

	pages := make([]string, page_count)
	ret_err := error(nil)

	create_page := func(loop_svg_group []string, loop_idx int) {
		wg.Add(1)
		go func(svg_group []string, idx int) {
			page, err := FitSvgsToPaper(format, svg_group)

			if err != nil {
				mut.Lock()
				ret_err = err
				mut.Unlock()
			} else {
				pages[idx] = page
			}

			wg.Done()
		}(loop_svg_group, loop_idx)
	}

	for loop_idx := 0; loop_idx < page_count-1; loop_idx++ {
		loop_svg_group := svgs[objects_on_page*loop_idx : objects_on_page*(loop_idx+1)]
		create_page(loop_svg_group, loop_idx)
	}

	create_page(svgs[objects_on_page*(page_count-1):], page_count-1)

	wg.Wait()

	if ret_err != nil {
		return []string{}, ret_err
	}

	return pages, nil
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

	papers, err := createPages(format, svgs)

	if err != nil {
		return []string{}, err
	}

	return papers, nil
}
