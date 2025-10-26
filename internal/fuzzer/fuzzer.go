package fuzzer

type FuzzResult struct {
	Path       string
	StatusCode int
}

func Fuzz(url, wordlist string) ([]FuzzResult, error) {
	// TODO: 实现模糊测试逻辑
	return []FuzzResult{
		{Path: "/admin", StatusCode: 403},
		{Path: "/admin.php", StatusCode: 404},
		{Path: "/robots.txt", StatusCode: 200},
	}, nil
}
