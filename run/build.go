package docs_blog_engine_run

func Build() error {
	return runCommand("./app/", "npm", "run", "build")
}
func Install() error {
	return runCommand("./app/", "npm", "install")
}
