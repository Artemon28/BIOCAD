package structures

type Device struct {
	Id        int
	Mqtt      string
	Invid     string
	UnitGuid  string
	MsgId     string
	Text      string
	Context   string
	Class     string
	Level     int
	Area      string
	Addr      string
	Block     string
	Type      string
	Bit       int
	InvertBit int
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type File struct {
	Id   int
	Name string
}
