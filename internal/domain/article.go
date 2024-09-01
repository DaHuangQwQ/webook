package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author

	ImgUrl string
	Type   string

	Status ArticleStatus
}

func (a Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) > 100 {
		return string(cs[:100])
	}
	return string(cs)
}

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

type ArticleStatus uint8

func (a ArticleStatus) ToUInt8() uint8 {
	return uint8(a)
}

type Author struct {
	Id   int64
	Name string
}
