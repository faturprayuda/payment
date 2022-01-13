package Log

type Log struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type LogHistory struct {
	Id         int    `json:"id"`
	Log        string `json:"log"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}