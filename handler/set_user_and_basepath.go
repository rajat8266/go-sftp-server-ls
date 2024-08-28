package handler

var (
	User     string
	BasePath string
)

func SetUserAndBasePath(user string, path string) {
	User = user
	BasePath = path
}
