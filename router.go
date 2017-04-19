package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter injects all the routes and returns a new router
func NewRouter() *mux.Router {

	router := mux.NewRouter()

	injectDeprecattedRoutes(router)

	injectRoutes(router)

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

func injectRoutes(router *mux.Router) {
	router.HandleFunc("/{artist}/tattoo", Tattoos).Methods("GET")
	router.HandleFunc("/{artist}/tattoo.toml", TattooToml).Methods("GET")
	router.HandleFunc("/{artist}/tattoo/{work_name}", CreateTattoo).Methods("POST")
	router.HandleFunc("/{artist}/tattoo/{work_name}", DeleteTattoo).Methods("DELETE")

	router.HandleFunc("/{artist}/henna", Hennas).Methods("GET")
	router.HandleFunc("/{artist}/henna.toml", HennaToml).Methods("GET")
	router.HandleFunc("/{artist}/henna/{work_name}", CreateHenna).Methods("POST")
	router.HandleFunc("/{artist}/henna/{work_name}", DeleteHenna).Methods("DELETE")

	router.HandleFunc("/{artist}/design", Designs).Methods("GET")
	router.HandleFunc("/{artist}/design.toml", DesignsToml).Methods("GET")
	router.HandleFunc("/{artist}/design/{work_name}", CreateDesign).Methods("POST")
	router.HandleFunc("/{artist}/design/{work_name}", DeleteDesign).Methods("DELETE")

	router.HandleFunc("/{artist}/piercing", Piercing).Methods("GET")
	router.HandleFunc("/{artist}/piercing.toml", PiercingToml).Methods("GET")
	router.HandleFunc("/{artist}/piercing/{work_name}", CreatePiercing).Methods("POST")
	router.HandleFunc("/{artist}/piercing/{work_name}", DeletePiercing).Methods("DELETE")

	router.HandleFunc("/{artist}/dreadlocks", Locks).Methods("GET")
	router.HandleFunc("/{artist}/dreadlocks.toml", LocksToml).Methods("GET")
	router.HandleFunc("/{artist}/dreadlocks/{work_name}", CreateDreadlocks).Methods("POST")
	router.HandleFunc("/{artist}/dreadlocks/{work_name}", DeleteDreadlocks).Methods("DELETE")

	for _, t := range []string{"tattoo", "henna", "piercing", "design", "dreadlocks", "art", "food"} {
		router.HandleFunc("/{name}/"+t, ArtistArtwork(t)).Methods("GET")
		router.HandleFunc("/{name}/"+t+"/refresh", ArtistArtworkRefresh(t)).Methods("GET")
	}

	router.HandleFunc("/{name}/all/refresh", ArtistArtworkRefreshAll).Methods("GET")

	router.HandleFunc("/artists", Artists).Methods("GET")

}
