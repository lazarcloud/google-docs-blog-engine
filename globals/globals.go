package globals

var (
	DescriptionKeyword = "DESCRIPTION"
	DefaultDescription = "No description."
	ImagesRoot         = "./app/public/images"
	DefaultImagePath   = "./app/public/blog-placeholder-about.jpg"
	BuildDir           = "./app/dist"
	StaticDir          = "./web"
)

func init() {
	if DescriptionKeyword[len(DescriptionKeyword)-1] != ' ' {
		DescriptionKeyword += " "
	}
}
