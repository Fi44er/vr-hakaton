package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

// Структура для передачи данных в HTML-шаблон
type EmailData struct {
	Name        string
	TeamName    string
	TeamID      string
	Email       string
	PhoneNumber string
	Role        string
}

// Функция для отправки письма с HTML-шаблоном
func Mailer(mails []string, fio string, teamName string, teamID string, email string, phoneNumber string, role string) {

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Адрес отправителя и пароль
	from := "fakeroot94@gmail.com"
	password := "memm cchg tscz gioi"

	// Данные для передачи в шаблон
	data := EmailData{
		Name:        fio, // Вы можете передавать имя пользователя динамически
		TeamName:    teamName,
		TeamID:      teamID,
		Email:       email,
		PhoneNumber: phoneNumber,
		Role:        role,
	}

	// Чтение HTML-шаблонов
	tmpl1, err := template.ParseFiles("pkg/mailer/hackaton.html")
	if err != nil {
		fmt.Println("Ошибка загрузки шаблона 1:", err)
		return
	}

	tmpl2, err := template.ParseFiles("pkg/mailer/hackaton2.html")
	if err != nil {
		fmt.Println("Ошибка загрузки шаблона 2:", err)
		return
	}

	// Отправка сообщений
	for i, to := range mails {
		var tmpl *template.Template
		if i == 0 {
			tmpl = tmpl1
		} else {
			tmpl = tmpl2
		}

		// Буфер для хранения результата обработки шаблона
		var body bytes.Buffer

		// Добавляем MIME заголовки
		body.Write([]byte("MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"))
		body.Write([]byte(fmt.Sprintf("Subject: Регистрация на Хакатон\r\n\r\n")))

		// Применяем шаблон к данным
		err = tmpl.Execute(&body, data)
		if err != nil {
			fmt.Println("Ошибка выполнения шаблона:", err)
			return
		}

		// Аутентификация
		auth := smtp.PlainAuth("", from, password, smtpHost)

		// Отправка сообщения
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, body.Bytes())
		if err != nil {
			fmt.Println("Ошибка отправки письма:", err)
			return
		}

		fmt.Println("Письмо успешно отправлено")
	}
}
