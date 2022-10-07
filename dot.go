package godots

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"
	"text/template"
)

var runtimeVars = map[string]string{
	"GOOS":   runtime.GOOS,
	"GOARCH": runtime.GOARCH,
}

func (dots Dots) All(cfg fs.FS, vars *Variables, output string) error {
	for _, name := range dots.List() {
		if err := dots.Run(name, cfg, vars, output); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func (dots Dots) Run(name string, cfg fs.FS, vars *Variables, output string) error {
	for _, dot := range dots {
		if dot.Name != name {
			continue
		}
		return dot.Render(cfg, BaseWriter(output), vars)
	}
	return fmt.Errorf("%s is not defined", name)
}

func (dots Dots) List() (names []string) {
	for _, dot := range dots {
		names = append(names, dot.Name)
	}
	return names
}

type DotFileWriter func(string) (io.WriteCloser, error)

func BaseWriter(root string) func(dst string) (io.WriteCloser, error) {
	return func(dst string) (io.WriteCloser, error) {
		if err := os.MkdirAll(path.Dir(path.Join(root, dst)), 0o750); err != nil && !os.IsExist(err) {
			return nil, err
		}
		return os.OpenFile(path.Join(root, dst), os.O_CREATE|os.O_WRONLY, 0o640)
	}
}

func (dot Dot) Render(cfg fs.FS, df DotFileWriter, vars *Variables) error {
	tpl := template.Must(template.ParseFS(cfg, dot.Templates...))

	for tSrc, tDst := range dot.FileMap {
		dst := os.ExpandEnv(tDst)

		log.Printf("rendering %s => %s", tSrc, dst)
		fd, err := df(os.ExpandEnv(tDst))
		if err != nil {
			log.Print(err)
			continue
		}
		tpl.ExecuteTemplate(fd, tSrc, map[string]interface{}{"var": vars.Global, "runtime": runtimeVars})
		fd.Close()
	}
	return nil
}
