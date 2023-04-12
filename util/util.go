package util

//HandleError handles errors
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
