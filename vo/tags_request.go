package vo

type CreateTagRequest struct {
	TagName string `json:"tag_name" binding:"required"`
}