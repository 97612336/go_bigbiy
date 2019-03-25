package handlers

import (
	"go_bigbiy/handlers/go_bigbiy"
	"net/http"
)

func MyUrls() {
	http.HandleFunc("/", go_bigbiy.Index_page)
	http.HandleFunc("/detail",go_bigbiy.Aricle_detail)
}
