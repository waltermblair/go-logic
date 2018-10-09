package logic

type Message struct {
	Header 		int 		`json:"header"`
	Body 		MessageBody `json:"body"`
}

type MessageBody struct {
	Configs 	[]Config 	`json:"configuration"`
	Input 		[]bool 		`json:"input"`
}

type Config struct {
	ID 			int 		`json:"id"`
	Status 		string 		`json:"status"`
	Function 	string		`json:"function"`
	NextKeys 	[]int		`json:"next-keys"`
}