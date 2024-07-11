package sendEmail

import (
	orderDB "miniMarket/db/orderDB"
	"miniMarket/db/userDB"

	"gopkg.in/mail.v2"
	"gorm.io/gorm"
)

func sendEmail(to, subject, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUser := "i.palik.work@gmail.com"
	smtpPassword := "gspx aukj qkqu xrtg"

	m := mail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmailRegister(email string) {
	sendEmail(email, "Благодарим за регистрацию!", "Привет! Рады видеть тебя на нашем сайте.")

}

func SendEmailOrder(order orderDB.Order, db *gorm.DB) {
	var user userDB.User
	db.Where("id= ?", order.IdUser).First(&user)
	sendEmail(user.Email, "Ваш заказ принят!", "Ссылка на заказ: "+"http://localhost:8080/order/"+order.IdOrder)
}
