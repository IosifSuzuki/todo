package request

type TodoForm struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
