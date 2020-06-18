package echo

import (
	"mime"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/webx-top/com"
)

var workDir string

func SetWorkDir(dir string) {
	if len(dir) == 0 {
		if len(workDir) == 0 {
			setWorkDir()
		}
		return
	}
	if !strings.HasSuffix(dir, FilePathSeparator) {
		dir += FilePathSeparator
	}
	workDir = dir
}

func setWorkDir() {
	workDir, _ = os.Getwd()
	workDir = workDir + FilePathSeparator
}

func init() {
	if len(workDir) == 0 {
		setWorkDir()
	}
}

func Wd() string {
	if len(workDir) == 0 {
		setWorkDir()
	}
	return workDir
}

// HandlerName returns the handler name
func HandlerName(h interface{}) string {
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}
	return t.String()
}

// HandlerPath returns the handler path
func HandlerPath(h interface{}) string {
	v := reflect.ValueOf(h)
	t := v.Type()
	switch t.Kind() {
	case reflect.Func:
		return runtime.FuncForPC(v.Pointer()).Name()
	case reflect.Ptr:
		t = t.Elem()
		fallthrough
	case reflect.Struct:
		return t.PkgPath() + `.` + t.Name()
	}
	return ``
}

func HandlerTmpl(handlerPath string) string {
	name := path.Base(handlerPath)
	var r []string
	var u []rune
	for _, b := range name {
		switch b {
		case '*', '(', ')':
			continue
		case '-':
			goto END
		case '.':
			r = append(r, string(u))
			u = []rune{}
		default:
			u = append(u, b)
		}
	}

END:
	if len(u) > 0 {
		r = append(r, string(u))
		u = []rune{}
	}
	for i, s := range r {
		r[i] = com.SnakeCase(s)
	}
	return `/` + strings.Join(r, `/`)
}

// Methods returns methods
func Methods() []string {
	return methods
}

// ContentTypeByExtension returns the MIME type associated with the file based on
// its extension. It returns `application/octet-stream` incase MIME type is not
// found.
func ContentTypeByExtension(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); len(t) == 0 {
		t = MIMEOctetStream
	}
	return
}

func static(r RouteRegister, prefix, root string) {
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	h := func(c Context) error {
		name := filepath.Join(root, c.Param("*"))
		if !strings.HasPrefix(name, root) {
			return ErrNotFound
		}
		return c.File(name)
	}
	if prefix == "/" {
		r.Get(prefix+"*", h)
	} else {
		r.Get(prefix+"/*", h)
	}
}
