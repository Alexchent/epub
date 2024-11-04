package epub

import (
	"errors"
	"io"
	"strings"
)

// ReadFile 返回一个文件的内容
func (p *Book) ReadFile(n string) (string, error) {
	part := strings.Split(n, "#")
	n = part[0]
	src := p.filename(n)
	fd, err := p.open(src)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	b, err := io.ReadAll(fd)
	if err != nil {
		return "", err
	}

	// 返回整个文件内容
	return string(b), nil
}

// ReadAll 返回所有章节的内容
func (p *Book) ReadAll() (string, error) {
	if p == nil {
		return "", errors.New("nil pointer receiver")
	}

	var content string
	readFiles := make(map[string]bool)
	var readAll func(points []NavPoint) error
	readAll = func(points []NavPoint) error {
		for _, point := range points {
			if point.Content.Src != "" {
				src := point.Content.Src
				if strings.Contains(src, "#") {
					parts := strings.Split(src, "#")
					if len(parts) != 2 {
						return errors.New("路径不止一个锚点:" + src)
					}
					src = parts[0]
				}
				if readFiles[src] {
					continue
				}
				readFiles[src] = true
				chapter, err := p.ReadFile(src)
				if err != nil {
					return err
				}
				content += chapter
			}
			if len(point.Points) > 0 {
				if err := readAll(point.Points); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := readAll(p.Ncx.Points); err != nil {
		return "", err
	}

	return content, nil
}
