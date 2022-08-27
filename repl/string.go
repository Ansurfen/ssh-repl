package repl

func Strip(input string) string {
	if input == "" {
		return ""
	}
	args := ""
	flag := true
	for _, ch := range input {
		if ch == ' ' && flag {
			args += string(ch)
			flag = false
		}
		if ch != ' ' {
			args += string(ch)
			flag = true
		}
	}
	return args
}