package list

func sendCountUpdate(channel chan int, value int) {
	select {
	case channel <- value:
		// Message sent successfully
	default:
		// Failed to send message. The buffer is probably full
	}
}
