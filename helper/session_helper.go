package helper

import "go-toko-kovan-al/config"

func GetSessionID(sessionID string, config config.Config) (uint, error) {
	var idUser uint

	idUser, err := ClaimTokenHelper(sessionID, config)

	// idUint, err := strconv.Atoi(sessionID)
	if err != nil {
		return 0, err
	}

	// idUser = uint(idUint)

	return idUser, nil
}

// func FormError(key string, value []ErrorResponse, status string, session *session.Session) {
// 	keyNew := fmt.Sprintf("%s-form", key)
// 	stu := fmt.Sprintf("%s-status", key)

// 	if session.Get(key) != nil {
// 		session.Delete(key)
// 	}

// 	session.Set(keyNew, value)
// 	session.Set(stu, status)
// 	session.SetExpiry(time.Second * 3)
// 	session.Save()
// }

// func AlertMassage(key, msg, status string, session *session.Session) {

// 	defer session.Save()

// 	session.Destroy()

// 	stu := fmt.Sprintf("%s-status", key)
// 	// timeExpaiyed := fmt.Sprintf("%s-time", key)
// 	session.Set(key, msg)
// 	session.Set(stu, status)
// 	// session.Set(timeExpaiyed, time.Now())
// 	session.SetExpiry(time.Second * 2)

// }

// func AlertSessionDestrory(key string, session *session.Session) error {

// 	dataSessionTime := session.Get(key)
// 	fmt.Println(dataSessionTime)
// 	time.Parse("02/01/2006", dataSessionTime.(string))

// 	return nil

// }

// func FormDeleteError(key string, session *session.Session) {
// 	defer session.Save()
// 	err := session.Destroy()
// 	if err != nil {
// 		fmt.Println("form session destroy", err.Error())
// 		session.Delete(key)
// 	}

// }

// func AlertMassageError(key, value, status string) fiber.Handler {

// 	return func(c *fiber.Ctx) error {
// 		data[]
// 	cookie := new(fiber.Cookie)
// 	cookie.Name = "alert-Value"
// 	cookie.Value = idEnkripsiString
// 	cookie.Expires = time.Now().Add(24 * time.Hour)
// 	cookie.HTTPOnly = true
// 	cookie.Secure = true
// helper.AlertMassage("msg-alert-new-user", "User berhasil login!.", "success", session)

// 	c.Cookie(cookie)
// return
// 	}
// }
