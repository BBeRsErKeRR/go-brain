package errors

type AnsibleInventoryException struct {
	Message string
}

func (e *AnsibleInventoryException) Error() string {
	return e.Message
}
