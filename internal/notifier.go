package internal

const (
	UpdateConfig     = "update-config"
	UpdateCategories = "update-categories"
)

type WebsocketNotifier struct {
	channels []chan string
}

func NewWebsocketNotifier() *WebsocketNotifier {
	return &WebsocketNotifier{channels: make([]chan string, 0)}
}

func (n *WebsocketNotifier) NewChannel() chan string {
	ch := make(chan string)
	n.channels = append(n.channels, ch)
	return ch
}

func (n *WebsocketNotifier) RemoveChannel(ch chan string) {
	for i, c := range n.channels {
		if c == ch {
			close(ch)
			n.channels = append(n.channels[:i], n.channels[i+1:]...)
		}
	}
}

func (n *WebsocketNotifier) Update(message string) {
	for _, c := range n.channels {
		c <- message
	}
}
