package data

type SuccessData struct {
	Success bool `json:"success"`
}

func NewSuccessData() SuccessData {
	return SuccessData{
		Success: true,
	}
}
