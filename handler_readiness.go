package main

import "net/http"

func (apiCfg *apiConfig) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
