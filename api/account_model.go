package api

type (
	createAccountRequest struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required,oneof=USD EUR"`
	}

	getAccountRequest struct {
		ID int64 `json:"id" binding:"required,min=1"`
	}

	listAccountRequest struct {
		Paging listPaging `json:"paging"`
	}

	listPaging struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}
)
