package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

type MLTrainingService struct {
	db *gorm.DB
	// 注入其他依赖服务
	dataService  *DataProcessingService
	modelService *MLModelService
}

// 创建训练任务
func (s *MLTrainingService) CreateTrainingJob(datasetID uint, modelType string, algorithm string, params map[string]interface{}) (*model.MLTrainingJob, error) {
	// 验证数据集
	var dataset model.MLDataset
	if err := s.db.First(&dataset, datasetID).Error; err != nil {
		return nil, err
	}

	// 验证参数
	if err := s.validateTrainingParams(modelType, algorithm, params); err != nil {
		return nil, err
	}

	// 创建训练任务
	paramsJSON, _ := json.Marshal(params)
	job := &model.MLTrainingJob{
		DatasetID:  datasetID,
		ModelType:  modelType,
		Algorithm:  algorithm,
		Parameters: string(paramsJSON),
		Status:     "pending",
		StartTime:  time.Now(),
	}

	if err := s.db.Create(job).Error; err != nil {
		return nil, err
	}

	// 异步启动训练
	go s.runTraining(job)

	return job, nil
}

// 执行训练
func (s *MLTrainingService) runTraining(job *model.MLTrainingJob) {
	job.Status = "running"
	s.db.Save(job)

	// 加载数据集
	dataset, err := s.dataService.LoadDataset(job.DatasetID)
	if err != nil {
		s.handleTrainingError(job, err)
		return
	}

	// 数据预处理
	features, labels, err := s.dataService.PreprocessData(dataset)
	if err != nil {
		s.handleTrainingError(job, err)
		return
	}

	// 训练模型
	model, metrics, err := s.modelService.TrainModel(job.ModelType, job.Algorithm, features, labels, job.Parameters)
	if err != nil {
		s.handleTrainingError(job, err)
		return
	}

	// 保存训练结果
	metricsJSON, _ := json.Marshal(metrics)
	job.Status = "completed"
	job.EndTime = time.Now()
	job.Metrics = string(metricsJSON)
	job.Progress = 100
	s.db.Save(job)

	// 创建模型评估
	s.createModelEvaluation(job, model, metrics)
}

// 创建模型评估
func (s *MLTrainingService) createModelEvaluation(job *model.MLTrainingJob, models interface{}, metrics map[string]interface{}) error {
	// 获取测试数据
	testData, err := s.dataService.GetTestData(job.DatasetID)
	if err != nil {
		return err
	}

	// 评估模型
	evalMetrics, err := s.modelService.EvaluateModel(models, testData)
	if err != nil {
		return err
	}

	// 计算特征重要性
	featureImportance, err := s.modelService.CalculateFeatureImportance(models)
	if err != nil {
		return err
	}

	// 创建评估记录
	evaluation := &model.MLEvaluation{
		ModelID:     job.ID,
		DatasetID:   job.DatasetID,
		Metrics:     string(evalMetrics),
		Features:    string(featureImportance),
		EvaluatedAt: time.Now(),
	}

	return s.db.Create(evaluation).Error
}

// 获取训练状态
func (s *MLTrainingService) GetTrainingStatus(jobID uint) (map[string]interface{}, error) {
	var job model.MLTrainingJob
	if err := s.db.First(&job, jobID).Error; err != nil {
		return nil, err
	}

	var metrics map[string]interface{}
	if job.Metrics != "" {
		json.Unmarshal([]byte(job.Metrics), &metrics)
	}

	return map[string]interface{}{
		"status":   job.Status,
		"progress": job.Progress,
		"metrics":  metrics,
		"error":    job.Error,
	}, nil
}

// 获取训练历史
func (s *MLTrainingService) GetTrainingHistory(modelType string, limit int) ([]map[string]interface{}, error) {
	var jobs []model.MLTrainingJob
	if err := s.db.Where("model_type = ?", modelType).
		Order("created_at DESC").
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	var history []map[string]interface{}
	for _, job := range jobs {
		var metrics map[string]interface{}
		json.Unmarshal([]byte(job.Metrics), &metrics)

		history = append(history, map[string]interface{}{
			"id":         job.ID,
			"dataset_id": job.DatasetID,
			"algorithm":  job.Algorithm,
			"status":     job.Status,
			"metrics":    metrics,
			"start_time": job.StartTime,
			"end_time":   job.EndTime,
		})
	}

	return history, nil
}

// 验证训练参数
func (s *MLTrainingService) validateTrainingParams(modelType string, algorithm string, params map[string]interface{}) error {
	// TODO: 实现参数验证逻辑
	return nil
}

// 处理训练错误
func (s *MLTrainingService) handleTrainingError(job *model.MLTrainingJob, err error) {
	job.Status = "failed"
	job.Error = err.Error()
	job.EndTime = time.Now()
	s.db.Save(job)
}
