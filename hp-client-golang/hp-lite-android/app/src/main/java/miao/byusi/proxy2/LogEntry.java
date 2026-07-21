package miao.byusi.proxy2;

import java.util.Date;

public class LogEntry {
    private Date timestamp;
    private String content;

    public LogEntry(Date timestamp, String content) {
        this.timestamp = timestamp;
        this.content = content;
    }

    public Date getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(Date timestamp) {
        this.timestamp = timestamp;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }
}