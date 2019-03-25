package handlers

import (
	"bigbiy_web/handlers/bigbiy_www"
	"net/http"
)

func MyUrls() {
	//http.HandleFunc("/", bigbiy_www.Show_all_message)
	//http.HandleFunc("/novel_v2", bigbiy_www.Nvl_v2)
	//http.HandleFunc("/novel_detail_v2", bigbiy_www.Nvl_detail)
	//http.HandleFunc("/chapter_detail_v2", bigbiy_www.Chapter_detail)
	//http.HandleFunc("/search_love_nvl_v2",bigbiy_www.Search_love_nvl)
	http.HandleFunc("/", bigbiy_www.New_index_page)
	http.HandleFunc("/detail",bigbiy_www.Aricle_detail)
}
