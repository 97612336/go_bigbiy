package bigbiy_www

import (
	"bigbiy_web/config"
	"bigbiy_web/models"
	"bigbiy_web/util"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

//搜索喜欢的小说
func Search_love_nvl(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	var data = make(map[string]interface{})
	search_words := util.Get_argument(r, "search_words")
	fmt.Println(search_words)
	template_path := config.Template_path + "search_nvl_v2.html"
	util.Render_template(w, template_path, data)
}

//去往章节详情页的接口
func Chapter_detail(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	var data = make(map[string]interface{})
	book_id := util.Get_argument(r, "book_id")
	chapter_id := util.Get_argument(r, "chapter_id")
	one_banner := Get_banner_by_id(util.String_to_int(book_id))
	book_name := one_banner.Name
	//获取redis中的热搜词
	hot_words := util.Get_redis("hot_words")
	data["hot_words"] = book_name + "," + hot_words
	sql_str := "select name,chapter_text from chapter where id=?;"
	rows, err := util.DB.Query(sql_str, chapter_id)
	util.CheckErr(err)
	var text_str string
	var name string
	for rows.Next() {
		rows.Scan(&name, &text_str)
	}
	var text_list []string
	util.Json_to_object(text_str, &text_list)
	data["text"] = text_list
	data["chapter_name"] = name
	//判断是否有上一章
	has_per, per_id := Has_next_or_pervious_chapter(chapter_id, book_id, 0)
	//判断是否有下一章
	has_next, next_id := Has_next_or_pervious_chapter(chapter_id, book_id, 1)
	data["has_per"] = has_per
	data["per_id"] = per_id
	data["has_next"] = has_next
	data["next_id"] = next_id
	data["book_id"] = book_id
	template_path := config.Template_path + "chapter_detail_v2.html"
	util.Render_template(w, template_path, data)
}

//去往书本详情页的接口
func Nvl_detail(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	var data = make(map[string]interface{})
	book_id := util.Get_argument(r, "book_id")
	//获取书名
	sql_str := "select id,name,book_img from book where id =?;"
	rows, err := util.DB.Query(sql_str, book_id)
	util.CheckErr(err)
	var book_info models.Banner_novel
	for rows.Next() {
		rows.Scan(&book_info.Book_id, &book_info.Name, &book_info.Img)
		Get_desc_by_book_id(&book_info)
	}
	data["book_info"] = book_info
	//获取该书下的所有章节
	sql_str2 := "select id,name from chapter where book_id=?;"
	rows2, err2 := util.DB.Query(sql_str2, book_id)
	util.CheckErr(err2)
	var chapters []models.Chapter_name
	for rows2.Next() {
		var one_chapter_name models.Chapter_name
		rows2.Scan(&one_chapter_name.Id, &one_chapter_name.Name)
		chapters = append(chapters, one_chapter_name)
	}
	data["chapters"] = chapters
	//获取redis中的图片
	img := util.Get_redis("biying_img")
	data["img"] = img
	//获取redis中的热搜词
	hot_words := util.Get_redis("hot_words")
	data["hot_words"] = book_info.Name + "," + hot_words
	template_path := config.Template_path + "nvl_detail_v2.html"
	util.Render_template(w, template_path, data)
}

//去往主页的方法
func Index_v2(w http.ResponseWriter, r *http.Request) {
	var data = make(map[string]interface{})
	//获取redis中的图片
	img := util.Get_redis("biying_img")
	data["img"] = img
	template_path := config.Template_path + "index_v2.html"
	util.Render_template(w, template_path, data)
}

//去往nvl的方法
func Nvl_v2(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	var data = make(map[string]interface{})
	last_book_id := util.Get_argument(r, "last_book_id", "0")
	first_book_id := util.Get_argument(r, "first_book_id", "0")
	kind := util.Get_argument(r, "kind", "1")
	search_words := util.Get_argument(r, "search_words", "0")
	//获取redis中的图片
	img := util.Get_redis("biying_img")
	data["img"] = img
	//获取redis中的热搜词
	hot_words := util.Get_redis("hot_words")
	data["hot_words"] = hot_words
	// 获取推荐书本
	banners := Get_banner()
	data["banner"] = banners
	//查询分页展示的书籍
	var books []models.Banner_novel
	var sql_str string
	var rows *sql.Rows
	var err error
	new_search_words := "'%" + search_words + "%'"
	//等于１表示下一页,也是默认的
	if kind == "1" {
		if search_words == "0" {
			sql_str = "select id,name,book_img,author from book where id >? and has_chapter=1 limit 15;"
			rows, err = util.DB.Query(sql_str, last_book_id)
		} else {
			sql_str = "select id,name,book_img,author from book where id >" + last_book_id +
				" and has_chapter=1 and name like " + new_search_words + " limit 15;"
			rows, err = util.DB.Query(sql_str)
		}
		util.CheckErr(err)
		var one_banner models.Banner_novel
		i := 1
		for rows.Next() {
			rows.Scan(&one_banner.Book_id, &one_banner.Name, &one_banner.Img, &one_banner.Author)
			if i == 1 {
				data["first_book_id"] = one_banner.Book_id
			}
			i = i + 1
			Get_desc_by_book_id(&one_banner)
			last_book_id := one_banner.Book_id
			data["last_book_id"] = last_book_id
			books = append(books, one_banner)
		}
	} else {
		//否则就是上一页
		if search_words == "0" {
			sql_str = "select id,name,book_img,author from book where id <? and has_chapter=1 order by id desc limit 15;"
			rows, err = util.DB.Query(sql_str, first_book_id)
		} else {
			sql_str = "select id,name,book_img,author from book where id <" + first_book_id +
				" and has_chapter=1 and name like " + new_search_words + " order by id desc limit 15;"
			rows, err = util.DB.Query(sql_str)
		}
		util.CheckErr(err)
		var one_banner models.Banner_novel
		i := 1
		for rows.Next() {
			rows.Scan(&one_banner.Book_id, &one_banner.Name, &one_banner.Img, &one_banner.Author)
			if i == 1 {
				data["last_book_id"] = one_banner.Book_id
			}
			i = i + 1
			Get_desc_by_book_id(&one_banner)
			last_book_id := one_banner.Book_id
			data["first_book_id"] = last_book_id
			books = append(books, one_banner)
		}
		//反转书本数组
		for i, j := 0, len(books)-1; i < j; i, j = i+1, j-1 {
			books[i], books[j] = books[j], books[i]
		}
	}
	data["books"] = books
	data["search_words"] = search_words
	//判断是否有下一页
	if len(books) < 15 {
		data["has_next"] = 0
	} else {
		data["has_next"] = 1
	}
	//判断是否有上一页
	data["has_per"] = Has_next_or_per_page_book(books, search_words, 0)
	//判断是否有下一页
	data["has_next"] = Has_next_or_per_page_book(books, search_words, 1)
	template_path := config.Template_path + "nvl_v2.html"
	util.Render_template(w, template_path, data)
}

func Get_desc_by_book_id(one_hot_novel *models.Banner_novel) {
	book_id := one_hot_novel.Book_id
	sql_str := "select chapter_text from chapter where book_id=" + strconv.Itoa(book_id) + " limit 1;"
	rows, err := util.DB.Query(sql_str)
	defer rows.Close()
	util.CheckErr(err)
	var text string
	for rows.Next() {
		rows.Scan(&text)
	}
	var text_list []string
	util.Json_to_object(text, &text_list)
	var desc string
	var i = 0
	for _, sentence := range text_list {
		desc = desc + sentence
		i = i + 1
		if i > 2 {
			break
		}
	}
	if util.String_length(desc) > 300 {
		desc = util.Splite_string(desc, 300)
	}
	one_hot_novel.Desc = desc + "......"
}

// 根据书本id查询数据库中的书本信息
func Get_banner_by_id(novel_id int) models.Banner_novel {
	sql_str := "select id,name,book_img,author from book where id=" + strconv.Itoa(novel_id) + ";"
	rows, err := util.DB.Query(sql_str)
	defer rows.Close()
	util.CheckErr(err)
	var one_banner models.Banner_novel
	for rows.Next() {
		rows.Scan(&one_banner.Book_id, &one_banner.Name, &one_banner.Img, &one_banner.Author)
		Get_desc_by_book_id(&one_banner)
	}
	return one_banner
}

//获取banner数据的方法
func Get_banner() []models.Banner_novel {
	banner_id_list := util.Get_banner_novel_id()
	var banners []models.Banner_novel
	for _, novel_id := range banner_id_list {
		one_banner := Get_banner_by_id(novel_id)
		banners = append(banners, one_banner)
	}
	return banners
}

//判断当前章节是否有上一章或者下一章
func Has_next_or_pervious_chapter(current_chapter_id string, book_id string, kind int) (int, int) {
	var num int
	var id int
	if kind == 1 {
		sql_str := "select count(1),id from chapter where id >? and book_id=? limit 1;"
		rows, err := util.DB.Query(sql_str, current_chapter_id, book_id)
		util.CheckErr(err)
		for rows.Next() {
			rows.Scan(&num, &id)
		}

	} else {
		sql_str := "select count(1) from chapter where id<?  and book_id=? limit 1;"
		rows, err := util.DB.Query(sql_str, current_chapter_id, book_id)
		util.CheckErr(err)
		for rows.Next() {
			rows.Scan(&num)
		}
		sql_str2 := "select id from chapter where id<? and book_id=? order by id DESC limit 1;"
		rows2, err2 := util.DB.Query(sql_str2, current_chapter_id, book_id)
		util.CheckErr(err2)
		for rows2.Next() {
			rows2.Scan(&id)
		}
	}
	if num > 0 {
		return 1, id
	} else {
		return 0, id
	}
}

//判断搜索结果中是否有上一页
func Has_next_or_per_page_book(books []models.Banner_novel, search_words string, kind int) int {
	var sql_str string
	new_search_words := "'%" + search_words + "%'"
	if len(books) < 1 {
		return 0
	}
	first_id := books[0].Book_id
	last_id := books[len(books)-1].Book_id
	if search_words == "0" {
		//kind等于１是查询是否有下一页
		if kind == 1 {
			sql_str = "select count(1) from book where id>" + util.Int_to_string(last_id) + " and has_chapter=1 limit 1;"
		} else {
			//	否则就是查询是否有上一页
			sql_str = "select count(1) from book where id<" + util.Int_to_string(first_id) + " and has_chapter=1 limit 1;"
		}
	} else {
		if kind == 1 {
			sql_str = "select count(1) from book where id>" + util.Int_to_string(last_id) +
				" and has_chapter=1 and name like " + new_search_words + " limit 1;"
		} else {
			sql_str = "select count(1) from book where id <" + util.Int_to_string(first_id) +
				" and has_chapter=1 and name like " + new_search_words + " limit 1;"
		}
	}
	rows, err := util.DB.Query(sql_str)
	util.CheckErr(err)
	var num int
	for rows.Next() {
		rows.Scan(&num)
	}
	if num > 0 {
		return 1
	} else {
		return 0
	}
}
