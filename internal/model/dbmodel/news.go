package dbmodel

//go:generate reform

//reform:news
type News struct {
	Id      int    `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}
