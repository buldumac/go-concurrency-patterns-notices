package generators

// 1. Generator of passed arguments values
func Generator(doneCh <-chan interface{}, values ...interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for _, value := range values {
			select {
			case <-doneCh: // exit when done channel is closed, ctx might be used too for that purpose
				return
			case stream <- value:
			}
		}
	}()
	return stream
}

// Take is useful when we want to receive a limited number of messages from some source channel
func Take(doneCh <-chan interface{}, provisionerStream <-chan interface{}, numOfItems int) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for i := 0; i < numOfItems; i++ {
			select {
			case <-doneCh:
				return
			case stream <- <-provisionerStream:
			}
		}
	}()
	return stream
}

// Repeat just loop through the values again and again and indefinitely push the values in the channel
func Repeat(doneCh <-chan interface{}, values ...interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			for _, value := range values {
				select {
				case <-doneCh:
					return
				case stream <- value:
				}
			}
		}
	}()
	return stream
}

// RepeatFunc write result of some function indefinitely into some stream
func RepeatFunc(doneCh <-chan interface{}, fn func() interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			select {
			case <-doneCh:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}
