package main

import (
	"net/http"
	"text/template"

	"github.com/howeyc/fsnotify"
)

func init() {
	// init cache
	go walkArchives()

	// watch for change
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}

		err = watcher.Watch("archives")
		if err != nil {
			panic(err)
		}
		defer watcher.Close()

		// FIXME 此 goroutine 不会自动推出，资源泄漏
		for {
			select {
			case <-watcher.Event:
				walkArchives()
			case err := <-watcher.Error:
				panic(err)
			}
		}
	}()
}

func main() {
	// static
	http.Handle("/image/", http.StripPrefix("/image", http.FileServer(http.Dir("static/image"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("static/css"))))
	http.Handle("/archives/", http.StripPrefix("/archives", http.FileServer(http.Dir("archives"))))

	// handler
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/atom.xml", feedHandler)

	panic(http.ListenAndServe(":80", nil))
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	render(w, "static/template/index.html", getAllArticles())
}

func aboutHandler(w http.ResponseWriter, _ *http.Request) {
	render(w, "static/template/about.html", nil)
}

func feedHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	tpl := template.Must(template.ParseFiles("static/template/atom.xml"))
	tpl.Execute(w, getAllArticles())
}
