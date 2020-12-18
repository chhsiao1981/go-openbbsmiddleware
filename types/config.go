package types

func config() {
	SERVICE_MODE = ServiceMode(setStringConfig("SERVICE_MODE", string(SERVICE_MODE)))

	HTTP_HOST = setStringConfig("HTTP_HOST", HTTP_HOST)

	URL_PREFIX = setStringConfig("URL_PREFIX", URL_PREFIX)

	JWT_SECRET = setBytesConfig("JWT_SECRET", JWT_SECRET)
	JWT_ISSUER = setStringConfig("JWT_ISSUER", JWT_ISSUER)
}
