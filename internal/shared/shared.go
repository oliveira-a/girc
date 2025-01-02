package shared

type CommandType int

const (
	Connect CommandType = iota
	Message
)

type Command struct {
	CommandType CommandType `json:"commandType"`
	From        string      `json:"from"`
	Content     string      `json:"content"`
}

type ClientMessage struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
