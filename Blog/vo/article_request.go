package vo

type CreateArticleRequest struct {
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
