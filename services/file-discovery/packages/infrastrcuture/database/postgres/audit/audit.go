package audit

type Operation string

const (
	CreateOperation  Operation = "C"
	ReadOperation 	 Operation = "R"
	UpdatedOperation Operation = "U"
	DeleteOperation  Operation = "D"
)
