package com.security.common.constant;

public class TaskConstants {
    // 任务状态
    public static final int TASK_STATUS_PENDING = 0;    // 待执行
    public static final int TASK_STATUS_RUNNING = 1;    // 执行中
    public static final int TASK_STATUS_FINISHED = 2;   // 已完成
    public static final int TASK_STATUS_FAILED = 3;     // 执行失败
    
    // 任务类型
    public static final int TASK_TYPE_ASSET = 1;        // 资产扫描
    public static final int TASK_TYPE_VULN = 2;         // 漏洞扫描
    
    // 漏洞等级
    public static final int VULN_LEVEL_LOW = 1;         // 低危
    public static final int VULN_LEVEL_MEDIUM = 2;      // 中危
    public static final int VULN_LEVEL_HIGH = 3;        // 高危
} 