package processor

type ThumbnailInfo struct {
	PageCount int
	Title     string
	Author    string
	Subject   string
}

func ExtractThumbnailInfo(pageCount int) ThumbnailInfo {
	return ThumbnailInfo{
		PageCount: pageCount,
		Title:     "",
		Author:    "",
		Subject:   "",
	}
}
