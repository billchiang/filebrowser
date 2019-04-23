package http

import (
	"net/http"
	"os"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/search"
)

var searchReqHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// response := []interface{}{}
	query := r.URL.Query().Get("query")
	var filesResponse []*files.FileInfo

	err := search.Search(d.user.Fs, r.URL.Path, query, d, func(path string, f os.FileInfo) error {

		//search result using NewFileInfo
		file, _ := files.NewFileInfo(files.FileOptions{
			Fs:      d.user.Fs,
			Path:    "/" + f.Name(),
			Modify:  false,
			Expand:  true,
			Checker: d,
		})
		filesResponse = append(filesResponse, file)

		return nil
	})
	//gen a dir from root to
	dirab, _ := files.NewFileInfo(files.FileOptions{
		Fs:      d.user.Fs,
		Path:    "/",
		Modify:  false,
		Expand:  true,
		Checker: d,
	})
	dirab.Listing.Items = filesResponse
	dirab.Listing.Sorting = d.user.Sorting
	dirab.Listing.ApplySort()

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return renderJSON(w, r, dirab)
})
