package globals

var (
	DescriptionKeyword = "DESCRIPTION"
	DefaultDescription = "No description."
	ImagesRoot         = "./app/public/images"
	DefaultImagePath   = "./app/public/blog-placeholder-about.jpg"
	BuildDir           = "./app/dist"
	StaticDir          = "./web"
	BackupDir = "./backup"
)

func init() {
	if DescriptionKeyword[len(DescriptionKeyword)-1] != ' ' {
		DescriptionKeyword += " "
	}
}
