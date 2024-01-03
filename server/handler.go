package server

import (
	"net/http"
)

type HttpHandler struct{}

func (HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//status_code := "200"
	//defer func(t time.Time) {
	//	logStr := fmt.Sprintf("[Access] %s | %s | %s | %s | %s |%s",
	//		time.Now().Format("2006/01/02 - 15:04:05"),
	//		//res.Header(),
	//		time.Since(t).String(),
	//		req.RemoteAddr,
	//		req.Method,
	//		status_code,
	//		req.RequestURI,
	//	)
	//}(time.Now())
	//defer func() {
	//	if err := recover(); err != nil {
	//		status_code = "500"
	//		res.WriteHeader(500)
	//		print(err)
	//		buff := debug.Stack()
	//
	//	}
	//}()

	//http.DefaultServeMux.ServeHTTP(res, req)
	mux.ServeHTTP(res, req)
}
