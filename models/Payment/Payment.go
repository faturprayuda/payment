package Payment

type PaymentTf struct {
	Id         int    `json:"id"`
	User_id    int    `json:"user_id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Nominal    int    `json:"Nominal"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}