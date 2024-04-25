package api

type (
	createAccountRequest struct {
		Currency string `json:"currency" binding:"required,currency"`
	}

	getAccountRequest struct {
		ID int64 `json:"id" binding:"required,min=1"`
	}

	listAccountRequest struct {
		Owner  string     `json:"owner"`
		Paging listPaging `json:"paging"`
	}

	listPaging struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}
)
