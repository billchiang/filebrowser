package http

import (
	"net/http"
	"os"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/search"
)

var searchReqHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	query := r.URL.Query().Get("query")
	var items []*files.FileInfo

	err := search.Search(d.user.Fs, r.URL.Path, query, d, func(path string, f os.FileInfo) error {

		items = append(items, &files.FileInfo{
			Fs:      d.user.Fs,
			Path:    r.URL.Path,
			Name:    path,
			ModTime: f.ModTime(),
			Mode:    f.Mode(),
			IsDir:   f.IsDir(),
			Size:    f.Size(),
		})

		return nil
	})
	if len(items) > 0 {
		searchReq, _ := files.NewFileInfo(files.FileOptions{
			Fs:      d.user.Fs,
			Path:    r.URL.Path,
			Modify:  false,
			Expand:  true,
			Checker: d,
		})

		searchReq.Listing.Items = items
		searchReq.Listing.Sorting = d.user.Sorting
		searchReq.NumDirs = 1
		searchReq.NumFiles = len(items)
		searchReq.Listing.ApplySort()
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return renderJSON(w, r, searchReq)
	}
	empty := []map[string]interface{}{}
	return renderJSON(w, r, empty)
})
