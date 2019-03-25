package bigbiy_www

import (
	"bigbiy_web/config"
	"bigbiy_web/models"
	"bigbiy_web/util"
	"math"
	"net/http"
	"strings"
)

func New_index_page(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	if r.Method == "GET" {
		var data = make(map[string]interface{})
		// 获取页数
		n := util.Get_argument(r, "n", "1")
		page := util.String_to_int(n)
		var page_size = 12
		start_num_str := util.Int_to_string((page - 1) * page_size)
		end_num_str := util.Int_to_string(page_size)
		// 从数据库读取数据
		sql_str := "select id,hot_word,title,info,imgs from articles order by id desc limit ?,?;"
		rows, err := util.DB.Query(sql_str, start_num_str, end_num_str)
		util.CheckErr(err)
		defer rows.Close()
		var articles []models.Article
		//遍历数据体
		for rows.Next() {
			//定义文章实体类
			var one_article models.Article
			// 定义图片字符串列表
			var imgs string
			err := rows.Scan(&one_article.Id, &one_article.Hot_word, &one_article.Title, &one_article.Info, &imgs)
			//进行判断，然后获取第一张图片
			if "[]" != imgs {
				var img_list []string
				new_imgs := strings.Replace(imgs, "'", "\"", -1)
				util.Json_to_object(new_imgs, &img_list)
				one_article.Img = img_list[0]
			} else {
				one_article.Img = "#"
			}
			util.CheckErr(err)
			articles = append(articles, one_article)
		}
		//获取页码数
		//总记录数
		count_num := Get_all_page_num()
		page_nums := Paginator(page, page_size, count_num)
		//进行数据传输，发送给HTML
		data["paginator"] = page_nums
		data["articles"] = articles
		data["current_page"] = page
		template_path := config.Template_path + "index.html"
		util.Render_template(w, template_path, data)
	}
}

func Aricle_detail(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 3)
	if r.Method == "GET" {
		var data = make(map[string]interface{})
		article_id_str := util.Get_argument(r, "article", "")
		//n := util.Get_argument(r, "n", "1")
		//page := util.String_to_int(n)
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
		// 把值赋予给data
		data["hot_word"] = hot_word
		data["title"] = title
		data["info"] = info
		var content_list []string
		new_content := strings.Replace(content, "'", "\"", -1)
		//打印输出
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
		//根据id推算是第几页
		//总记录数
		count_num := Get_all_page_num()
		page_size := 12
		//获取该记录是倒数第几个
		article_id := util.String_to_int(article_id_str)
		reci_num := count_num - article_id
		if reci_num <= 0 {
			data["page"] = 1
		} else {
			page := int(math.Ceil(float64(reci_num) / float64(page_size)))
			data["page"] = page
		}
		template_path := config.Template_path + "detail.html"
		util.Render_template(w, template_path, data)
	}
}
