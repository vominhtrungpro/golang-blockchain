package cli

import "log"

// handle all error
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
