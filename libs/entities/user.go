package entities

type User struct {
	PK           string `dynamodbav:"pk"`
	SK           string `dynamodbav:"sk"`
	Email        string `dynamodbav:"email"`
	Name         string `dynamodbav:"name"`
	Password     string `dynamodbav:"password"`
	PasswordSalt string `dynamodbav:"password_salt"`
}
