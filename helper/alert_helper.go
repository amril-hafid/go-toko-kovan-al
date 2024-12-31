package helper

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

func AlertString(massage interface{}, status interface{}) string {
	var msg string
	str := fmt.Sprintf("%v", massage)
	stu := fmt.Sprintf("%v", status)

	if stu == "success" {
		msg = fmt.Sprintf("<script> $(function() { var Toast = Swal.mixin({ toast: true, position: 'top-end', showConfirmButton: false, timer: 3000 });  $('.toastrDefaultSuccess').ready(function() { toastr.success('%s') }); }); </script>", str)

		// msg = fmt.Sprintf("<script> $(function () { var Toast = Swal.mixin({ toast: true, position: 'top-end', showConfirmButton: false, timer: 3000 });  $('.toastsDefaultSuccess').ready(function () {$(document).Toasts('create', {class: 'bg-success', title: 'Toast Title', subtitle: 'Subtitle', body: '%s' }) }); }); </script>", str)

	} else if stu == "error" {
		msg = fmt.Sprintf("<script> $(function() { var Toast = Swal.mixin({ toast: true, position: 'top-end', showConfirmButton: false, timer: 3000 });  $('.toastrDefaultError').ready(function() { toastr.error('%s') }); }); </script>", str)

		// msg = fmt.Sprintf("<script> $(function () { var Toast = Swal.mixin({ toast: true, position: 'top-end', showConfirmButton: false, timer: 3000 });  $('.toastsDefaultWarning').ready(function () {$(document).Toasts('create', {class: 'bg-warning', title: 'Toast Title', subtitle: 'Subtitle', body: '%s' }) }); }); </script>", str)

	} else {
		// msg = fmt.Sprintf("<script> $(function() { var Toast = Swal.mixin({ toast: true, position: 'top-end', showConfirmButton: false, timer: 10000 });  $('.toastrDefaultError').ready(function() { toastr.error('%s') }); }); </script>", str)

		msg = ""

	}
	return msg
}

func AlertMassage(key, msg, status string, session *session.Session) {

	defer session.Save()

	session.Destroy()

	stu := fmt.Sprintf("%s-status", key)
	// timeExpaiyed := fmt.Sprintf("%s-time", key)
	session.Set(key, msg)
	session.Set(stu, status)
	// session.Set(timeExpaiyed, time.Now())
	session.SetExpiry(time.Second * 2)
}

func FormDeleteError(key string, session *session.Session) {
	defer session.Save()
	err := session.Destroy()
	if err != nil {
		fmt.Println("form session destroy", err.Error())
		session.Delete(key)
	}

}

// <script> $(function() { var Toast = Swal.mixin({ toast: true, position: 'top-end', showConfirmButton: false, timer: 3000 }); document.getElementById('targetElement') $('.toastrDefaultSuccess').click(function() { toastr.success('Lorem ipsum dolor sit amet, consetetur sadipscing elitr.') }); }); </script>
