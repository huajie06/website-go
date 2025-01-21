package index

import "html/template"

// index page will show the most recent posts so the data will be in a slice
// what needed are the slice of each post/articles
// which will be each artical has a title and url
type Index_Content struct {
	Home_Content_HTML template.HTML
}
