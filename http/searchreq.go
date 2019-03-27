package http

import (
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/search"
	"log"
	"net/http"
	"os"
)

var searchReqHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	response := []map[string]interface{}{}
	query := r.URL.Query().Get("query")

	err := search.Search(d.user.Fs, r.URL.Path, query, d, func(path string, f os.FileInfo) error {

		file, _ := files.NewFileInfo(files.FileOptions{
			Fs:      d.user.Fs,
			Path:    "/" + f.Name(),
			Modify:  false,
			Expand:  true,
			Checker: d,
		})
		response = append(response, map[string]interface{}{
			"dir":       file.IsDir,
			"path":      file.Path,
			"extension": file.Extension,
		})
		log.Print(file)
		log.Print(file.Path)
		log.Print(file.Path)
		log.Print("This is search req class")
		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, response)
})
