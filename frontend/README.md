# Frontend 目录说明

## 📁 目录结构

```
frontend/
├── index.html      # 主HTML文件
├── app.js          # 前端逻辑
├── dist/           # 编译后的资源（Wails使用此目录）
│   ├── index.html  # 副本
│   └── app.js      # 副本
└── wailsjs/        # Wails自动生成的绑定
    ├── go/         # Go方法绑定
    └── runtime/    # 运行时API
```

## 🎯 开发流程

1. **编辑源文件**: 直接修改 `frontend/index.html` 和 `frontend/app.js`
2. **复制到dist**: 运行 `./sync-frontend.sh` 或手动复制文件到 `frontend/dist/`
3. **重新运行**: Wails 会自动检测变化并热重载

## 📝 注意事项

- Wails 会从 `frontend/dist/` 目录加载资源
- `index.html` 和 `app.js` 都需要存在于 `dist/` 目录
- 使用 `sync-frontend.sh` 脚本来同步文件

