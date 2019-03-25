package bigbiy_www

import (
	"bigbiy_web/config"
	"bigbiy_web/models"
	"bigbiy_web/util"
	"math"
	"net/http"
	"strings"
)

func Show_all_message(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	if r.Method == "GET" {
		var data = make(map[string]interface{})
		// 获取页数
		n := util.Get_argument(r, "n", "1")
		page := util.String_to_int(n)
		var page_size = 5
		start_num_str := util.Int_to_string((page - 1) * page_size)
		end_num_str := util.Int_to_string(page_size)
		// 从数据库读取数据
		sql_str := "select id,hot_word,title from articles order by id desc limit ?,?;"
		rows, err := util.DB.Query(sql_str, start_num_str, end_num_str)
		util.CheckErr(err)
		defer rows.Close()
		var articles []models.Article
		for rows.Next() {
			var one_article models.Article
			err := rows.Scan(&one_article.Id, &one_article.Hot_word, &one_article.Title)
			util.CheckErr(err)
			articles = append(articles, one_article)
		}
		count_num := Get_all_page_num()
		page_num := int(math.Ceil(float64(count_num) / float64(page_size)))
		var page_num_list []int
		for i := 1; i <= page_num; i++ {
			page_num_list = append(page_num_list, i)
		}
		data["page_num_list"] = page_num_list
		data["articles"] = articles
		data["current_page"] = page
		template_path := config.Template_path + "index.html"
		util.Render_template(w, template_path, data)
	}
}

// 得到文章总数量的方法
func Get_all_page_num() int {
	sql_str := "select count(1) from articles;"
	rows, err := util.DB.Query(sql_str)
	util.CheckErr(err)
	defer rows.Close()
	var num int
	for rows.Next() {
		err := rows.Scan(&num)
		util.CheckErr(err)
	}
	return num
}

// 点击进入文章详情页的方法
func Go_to_article_detail(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	if r.Method == "GET" {
		article_id_str := util.Get_argument(r, "id", "")
		n := util.Get_argument(r, "n", "1")
		page := util.String_to_int(n)
		sql_str := "select hot_word,title,info,content,imgs from articles where id=?;"
		rows, err := util.DB.Query(sql_str, article_id_str)
		util.CheckErr(err)
		defer rows.Close()
		var hot_word string
		var title string
		var info string
		var content string
		var imgs string
		for rows.Next() {
			err := rows.Scan(&hot_word, &title, &info, &content, &imgs)
			util.CheckErr(err)
		}
		var data = make(map[string]interface{})
		// 把值赋予给data
		data["hot_word"] = hot_word
		data["title"] = title
		data["info"] = info
		var content_list []string
		new_content := strings.Replace(content, "'", "\"", -1)
		util.Json_to_object(new_content, &content_list)
		data["content_list"] = content_list
		var img_list []string
		if "[]" != imgs {
			new_imgs := strings.Replace(imgs, "'", "\"", -1)
			util.Json_to_object(new_imgs, &img_list)
			data["img_list"] = img_list
		} else {
			data["img_list"] = []string{"#", "#"}
		}
		data["page"] = page
		template_path := config.Template_path + "detail.html"
		util.Render_template(w, template_path, data)
	}
}
