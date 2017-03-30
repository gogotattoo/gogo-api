package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter injects all the routes and returns a new router
func NewRouter() *mux.Router {

	router := mux.NewRouter()

	injectDeprecattedRoutes(router)

	router.HandleFunc("/upload", upload)
	// router.Handle("/uploaded/", http.StripPrefix("/uploaded/", http.FileServer(http.Dir("./upload/"))))
	router.PathPrefix("/uploaded/").Handler(http.StripPrefix("/uploaded/", http.FileServer(http.Dir("upload/"))))

	return router
}

func injectDeprecattedRoutes(router *mux.Router) {
	router.HandleFunc("/tattoo", Tattoos).Methods("GET")
	router.HandleFunc("/tattoo.toml", TattooToml).Methods("GET")
	router.HandleFunc("/tattoo/{id}", CreateTattoo).Methods("POST")
	router.HandleFunc("/tattoo/{id}", DeleteTattoo).Methods("DELETE")

	router.HandleFunc("/henna.toml", HennaToml).Methods("GET")
	router.HandleFunc("/henna", Hennas).Methods("GET")
	router.HandleFunc("/henna/{id}", CreateHenna).Methods("POST")
	router.HandleFunc("/henna/{id}", DeleteHenna).Methods("DELETE")

	router.HandleFunc("/design.toml", DesignsToml).Methods("GET")
	router.HandleFunc("/design", Designs).Methods("GET")
	router.HandleFunc("/design/{id}", CreateDesign).Methods("POST")
	router.HandleFunc("/design/{id}", DeleteDesign).Methods("DELETE")

	router.HandleFunc("/piercing.toml", PiercingToml).Methods("GET")
	router.HandleFunc("/piercing", Piercing).Methods("GET")
	router.HandleFunc("/piercing/{id}", CreatePiercing).Methods("POST")
	router.HandleFunc("/piercing/{id}", DeletePiercing).Methods("DELETE")

	router.HandleFunc("/dreadlocks.toml", LocksToml).Methods("GET")
	router.HandleFunc("/dreadlocks", Locks).Methods("GET")
	router.HandleFunc("/dreadlocks/{id}", CreateDreadlocks).Methods("POST")
	router.HandleFunc("/dreadlocks/{id}", DeleteDreadlocks).Methods("DELETE")

	for _, t := range []string{"tattoo", "henna", "piercing", "design", "dreadlocks"} {
		router.HandleFunc("/"+t+"/{name}", ArtistArtwork(t)).Methods("GET")
		router.HandleFunc("/"+t+"/{name}/refresh", ArtistArtworkRefresh(t)).Methods("GET")
	}

	router.HandleFunc("/all/{name}/refresh", ArtistArtworkRefreshAll).Methods("GET")

}
