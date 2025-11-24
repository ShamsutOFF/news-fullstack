package types

type Category struct {
	Name     string
	ImageURL string
}

type BannerCard struct {
	ImageURL    string
	Title       string
	Description string
}

type ArticleCard struct {
	ImageURL     string
	Title        string
	Description  string
	AuthorName   string
	AuthorAvatar string
	PublishDate  string
}
