package dbmodel

//go:generate reform

//reform:news_categories
type Categories struct {
	Id         int `reform:"id,pk"`
	NewsId     int `reform:"news_id"`
	CategoryId int `reform:"category_id"`
}
