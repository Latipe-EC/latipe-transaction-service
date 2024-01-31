package message

const (
	ORDER    = 1
	REGISTER = 2
	FORGOT   = 3
)

type EmailPurchaseMessage struct {
	EmailRecipient string `json:"email_recipient" validator:"email"`
	Name           string `json:"name"`
	OrderId        string `json:"order_id"`
}
