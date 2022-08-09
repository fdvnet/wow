package protocol

type Solution struct {
	Answer []byte
}

type Task struct {
	Difficulty int
	Nonce      []byte
}

type Message struct {
	Quote string
}
