package preprocessor

func Compiler(isCPP bool) (compiler, compilerFlag string) {
	compiler = "clang"
	compilerFlag = "-std=c99"
	if isCPP {
		compiler = "clang++"
		compilerFlag = "-std=c++98"
	}
	return
}
