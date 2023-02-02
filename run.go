package badgrlib

import "sync"

func createSvgs(format Format, data_source string) (chan string, chan error) {
	svgs := make(chan string)
	svgs_err := make(chan error)

	go func() {
		data := ParseTable(data_source)

		var wg sync.WaitGroup
		wg.Add(len(data.Data))

		for _, object_data := range data.Data {
			go func(object_data map[string]string) {
				svg, err := CreateSingleSvg(format, object_data)

				if err != nil {
					svgs_err <- err
				} else {
					svgs <- svg
				}

				wg.Done()
			}(object_data)
		}

		wg.Wait()

		close(svgs)
		close(svgs_err)
	}()

	return svgs, svgs_err
}

func createPapers(format Format, svgs chan string) (chan string, chan error) {
	papers := make(chan string)
	papers_err := make(chan error)

	objects_on_page := format.PaperFit.X * format.PaperFit.Y

	svg_groups := make(chan []string)
	go func() {
		svg_group := make([]string, 0, objects_on_page)

		for object_svg := range svgs {
			svg_group = append(svg_group, object_svg)

			if len(svg_group) == objects_on_page {
				svg_groups <- svg_group
				svg_group = make([]string, 0, objects_on_page)
			}
		}

		if len(svg_group) > 0 {
			svg_groups <- svg_group
		}

		close(svg_groups)
	}()

	go func() {
		var wg sync.WaitGroup

		for svg_group := range svg_groups {
			wg.Add(1)
			go func(svg_group []string) {
				paper, err := FitSvgsToPaper(format, svg_group)

				if err != nil {
					papers_err <- err
				} else {
					papers <- paper
				}

				wg.Done()
			}(svg_group)
		}

		wg.Wait()

		close(papers)
		close(papers_err)
	}()

	return papers, papers_err
}

func RunFormat(format_source string, data_source string) ([]string, error) {
	format, err := ParseFormat(format_source)

	if err != nil {
		return []string{}, err
	}

	svgs, svgs_err := createSvgs(format, data_source)

	papers, papers_err := createPapers(format, svgs)

	papers_slice := make([]string, 0)
	for paper := range papers {
		papers_slice = append(papers_slice, paper)
	}

	for svg_err := range svgs_err {
		return []string{}, svg_err
	}

	for paper_err := range papers_err {
		return []string{}, paper_err
	}

	return papers_slice, nil
}
