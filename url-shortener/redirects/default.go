package redirects

func defaultRedirects() Redirects {
	return map[string]string{
		"google": "http://google.com",
		"ya":     "http://yandex.ru",
		"yandex": "http://yandex.ru",
	}
}
