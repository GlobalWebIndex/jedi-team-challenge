package customerrors

type ResourceNotFoundErrorWrapper struct {
	OriginalError error
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err ResourceNotFoundErrorWrapper) Error() string {
	return ""
}

func (err ResourceNotFoundErrorWrapper) Unwrap() error {
	return err.OriginalError
}

type UserMismatchError struct {
	chatSessionID string
	userID        string
}

func NewUserMismatchError(chatSessionID, userID string) *UserMismatchError {
	return &UserMismatchError{
		chatSessionID: chatSessionID,
		userID:        userID,
	}
}

func (err UserMismatchError) Error() string {
	return "chatSession " + err.chatSessionID + " does not belong to user " + err.userID
}
