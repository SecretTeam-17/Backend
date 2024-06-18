package redirect

import (
	"fmt"
	"net/http"
)

func ToTLS(w http.ResponseWriter, r *http.Request) {
	url := "https://" + r.Host + r.RequestURI
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
