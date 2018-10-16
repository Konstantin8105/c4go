package preprocessor

func Compiler(isCPP bool) (compiler string, compilerFlag []string) {
	compiler = "clang"
	compilerFlag = []string{"-O0"}
	if isCPP {
		compiler = "clang++"
		compilerFlag = []string{"-std=c++98"}
	}
	return
}
