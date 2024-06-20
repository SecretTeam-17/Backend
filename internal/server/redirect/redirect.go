package redirect

import (
	"log"
	"net/http"
	"petsittersGameServer/internal/config"
)

// func ToTLS(w http.ResponseWriter, r *http.Request) {
// 	url := "https://" + r.Host + r.RequestURI
// 	fmt.Println(url)
// 	http.Redirect(w, r, url, http.StatusMovedPermanently)
// }

// ToTLS - перенаправляет запросы по http на https.
func ToTLS(cfg *config.Config) {
	rd := func(w http.ResponseWriter, r *http.Request) {
		url := "https://" + r.Host + r.RequestURI
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
	go func() {
		if err := http.ListenAndServe(cfg.Address, http.HandlerFunc(rd)); err != nil {
			log.Printf("failed to start redirect server: %s", err.Error())
		}
	}()
}
