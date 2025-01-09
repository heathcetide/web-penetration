package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"strings"
	"web_penetration/internal/model"
)

type UserImportExportService struct {
	db *gorm.DB
}

func NewUserImportExportService(db *gorm.DB) *UserImportExportService {
	return &UserImportExportService{db: db}
}

type UserExportData struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Role        string   `json:"role"`
	Status      int      `json:"status"`
	Groups      []string `json:"groups"`
	IsVerified  bool     `json:"is_verified"`
	CreatedAt   string   `json:"created_at"`
	LastLoginAt string   `json:"last_login_at"`
}

// 导出用户数据
func (s *UserImportExportService) ExportUsers(format string, userIDs []uint) ([]byte, error) {
	var users []model.User
	if err := s.db.Find(&users, userIDs).Error; err != nil {
		return nil, err
	}

	exportData := make([]UserExportData, 0)
	for _, user := range users {
		// 获取用户组
		var groups []string
		if err := s.db.Model(&user).
			Select("user_groups.name").
			Joins("JOIN user_group_members ON user_group_members.user_id = users.id").
			Joins("JOIN user_groups ON user_groups.id = user_group_members.group_id").
			Pluck("name", &groups).Error; err != nil {
			return nil, err
		}

		exportData = append(exportData, UserExportData{
			Username:    user.Username,
			Email:       user.Email,
			Phone:       user.Phone,
			Role:        user.Role,
			Status:      user.Status,
			Groups:      groups,
			IsVerified:  user.IsVerified,
			CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
			LastLoginAt: user.LastLoginAt.Format("2006-01-02 15:04:05"),
		})
	}

	switch format {
	case "json":
		return json.Marshal(exportData)
	case "csv":
		return s.exportToCSV(exportData)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func (s *UserImportExportService) exportToCSV(data []UserExportData) ([]byte, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{"Username", "Email", "Phone", "Role", "Status", "Groups", "IsVerified", "CreatedAt", "LastLoginAt"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, user := range data {
		record := []string{
			user.Username,
			user.Email,
			user.Phone,
			user.Role,
			fmt.Sprintf("%d", user.Status),
			strings.Join(user.Groups, "|"),
			fmt.Sprintf("%v", user.IsVerified),
			user.CreatedAt,
			user.LastLoginAt,
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return []byte(buf.String()), nil
}

// 导入用户数据
func (s *UserImportExportService) ImportUsers(reader io.Reader, format string) (int, error) {
	switch format {
	case "csv":
		return s.importFromCSV(reader)
	default:
		return 0, fmt.Errorf("unsupported format: %s", format)
	}
}

func (s *UserImportExportService) importFromCSV(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)

	// 跳过表头
	if _, err := csvReader.Read(); err != nil {
		return 0, err
	}

	imported := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return imported, err
		}

		user := model.User{
			Username: record[0],
			Email:    record[1],
			Phone:    record[2],
			Role:     record[3],
		}

		// 设置密码为默认值,后续需要用户修改
		user.Password = "123456"
		if err := user.HashPassword(); err != nil {
			return imported, err
		}

		if err := s.db.Create(&user).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue // 跳过重复数据
			}
			return imported, err
		}

		imported++
	}

	return imported, nil
}
