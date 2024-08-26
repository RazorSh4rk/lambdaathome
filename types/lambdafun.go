package types

type LambdaFun struct {
	// will be the endpoint
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Runtime string `json:"runtime"`
	Port    string `json:"port"`
	// ./data:/data,./folder:/app/folder
	Volume string `json:"volume"`
	Source string `json:"source"`
	ID     string `json:"id"`
}
