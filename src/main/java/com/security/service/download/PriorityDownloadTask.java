package com.security.service.download;

public class PriorityDownloadTask extends DownloadTask implements Comparable<PriorityDownloadTask> {
    private int priority; // 优先级：1-最高，5-最低
    private long createTime;
    
    @Override
    public int compareTo(PriorityDownloadTask other) {
        int priorityCompare = Integer.compare(this.priority, other.priority);
        if (priorityCompare != 0) {
            return priorityCompare;
        }
        return Long.compare(this.createTime, other.createTime);
    }

    public int getPriority() {
        return priority;
    }

    public void setPriority(int priority) {
        this.priority = priority;
    }

    public long getCreateTime() {
        return createTime;
    }

    public void setCreateTime(long createTime) {
        this.createTime = createTime;
    }
}