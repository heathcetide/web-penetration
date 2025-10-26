package filescan

type FileResult struct {
	Path   string
	Status string
}

func Scan(url string) ([]FileResult, error) {
	// TODO: 实现文件扫描逻辑
	return []FileResult{
		{Path: "/.env", Status: "404"},
		{Path: "/backup.sql", Status: "200"},
		{Path: "/config.php.bak", Status: "404"},
	}, nil
}
