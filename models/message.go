package models

//Message contains a body and whether it is a palindrome or not
type Message struct {
	ID           int    `json:"id"`
	Body         string `json:"body"`
	IsPalindrome bool   `gorm:"-" json:"isPalindrome"`
}
