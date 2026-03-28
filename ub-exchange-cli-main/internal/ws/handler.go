package ws

// Ws represents a WebSocket server that pushes real-time market data to clients.
type Ws interface {
	// Run starts the WebSocket server and begins broadcasting market data to connected clients.
	Run()
}
