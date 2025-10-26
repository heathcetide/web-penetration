package sqli

type ScanResult struct {
	Vulnerable bool
	Type       string
}

func Scan(url, parameter string) (*ScanResult, error) {
	// TODO: 实现SQL注入检测逻辑
	return &ScanResult{
		Vulnerable: false,
		Type:       "",
	}, nil
}
