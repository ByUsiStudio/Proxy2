package miao.byusi.proxy2;


import android.annotation.SuppressLint;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.os.Build;
import android.os.Handler;
import android.os.IBinder;
import android.os.Message;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.core.app.NotificationCompat;


import miao.byusi.proxy2.util.ConstConfig;
import miao.byusi.proxy2.util.SharedPreferencesUtil;

import java.text.SimpleDateFormat;
import java.util.Date;

import hp_android_lib.Callback;
import hp_android_lib.Hp_android_lib;


public class ProxyService extends Service {
    private static final int NOTIFICATION_ID = 1;
    private static final String channelId = "YOUR_CHANNEL_ID";
    private static final String channelName = "YOUR_CHANNEL_NAME";
    private NotificationManager notificationManager;
    private static boolean isStart = false;

    @Override
    public void onCreate() {
        Log.i("Kathy", "onCreate - Thread ID = " + Thread.currentThread().getId());
        super.onCreate();

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            // 当前设备运行的是 Android 8.0 或更高版本

            int importance = NotificationManager.IMPORTANCE_DEFAULT;
            NotificationChannel channel = new NotificationChannel(channelId, channelName, importance);
            notificationManager = getSystemService(NotificationManager.class);
            notificationManager.createNotificationChannel(channel);

            NotificationCompat.Builder builder = new NotificationCompat.Builder(this, channelId)
                    .setSmallIcon(R.drawable.ic_launcher)
                    .setContentTitle("Proxy2")
                    .setContentText("服务已经运行");
            startForeground(NOTIFICATION_ID, builder.build());
        } else {
            // 当前设备运行的是 Android 8.0 之前的版本
            NotificationCompat.Builder builder = new NotificationCompat.Builder(this, channelId)
                    .setContentTitle("Proxy2")
                    .setContentText("服务已经运行")
                    .setSmallIcon(R.drawable.ic_launcher)
                    .setPriority(NotificationCompat.PRIORITY_LOW);
            // 获取 NotificationManager 实例
            notificationManager = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
            startForeground(NOTIFICATION_ID, builder.build());
        }


    }


    private void updateNotification(String message) {
        // 更新通知的内容
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, channelId)
                .setContentTitle("Proxy2")
                .setContentText(message)
                .setSmallIcon(R.drawable.ic_launcher)
                .setPriority(NotificationCompat.PRIORITY_LOW);

        // 使用相同的通知 ID 更新通知
        notificationManager.notify(NOTIFICATION_ID, builder.build());
    }

    public Message getMessage(String msg) {
        Message message = new Message();
        message.what = -1;
        message.obj = msg;
        return message;
    }


    private StatusHandler handler = new StatusHandler();
    private Runnable runnableCode = new Runnable() {
        @Override
        public void run() {
            // 执行定时任务的操作
            // 在这里可以添加你想要执行的代码逻辑，例如发送通知、更新数据等
            new Thread() {
                @Override
                public void run() {
                    if (Hp_android_lib.getStatus()) {
                        handler.sendMessage(getMessage("true"));
                    } else {
                        handler.sendMessage(getMessage("false"));
                    }
                }
            }.start();
            // 每隔一定时间再次执行该任务
            handler.postDelayed(this, 1000 * 10); // 10秒
        }
    };

    @SuppressLint("HandlerLeak")
    public class StatusHandler extends Handler {

        private String getTime() {
            Date currentDate = new Date();
            // 创建 SimpleDateFormat 对象，指定要显示的时间格式
            SimpleDateFormat sdf = new SimpleDateFormat("HH:mm:ss");
            // 使用 SimpleDateFormat 格式化当前时间
            return sdf.format(currentDate);
        }

        @Override
        public void handleMessage(@NonNull Message msg) {
            if ("true".equals(msg.obj.toString())) {
                updateNotification("服务连接正常-检查时间:" + getTime());
            } else {
                updateNotification("服务连接异常-检查时间:" + getTime());
            }
        }
    }


    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        try {
            if (!isStart) {
                isStart = true;
                FirstFragment.RegHandler regHandler = new FirstFragment.RegHandler();
                String connect = SharedPreferencesUtil.getString(getApplicationContext(), ConstConfig.CONNECT, "");
                regHandler.sendMessage(getMessage("连接码：" + connect));
                new Thread(() -> {
                    Hp_android_lib.start(connect, new Callback() {
                        @Override
                        public void sendResult(String s) {
                            regHandler.sendMessage(getMessage(s));
                        }
                    });
                }).start();
                // 在服务启动时开始执行定时任务
                handler.postDelayed(runnableCode, 1000 * 10); // 10秒
            }
        } catch (Throwable e) {
        }
        return super.onStartCommand(intent, flags, startId);
    }

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }


    @Override
    public void onDestroy() {
        isStart = false;
        try {
            Hp_android_lib.close();
        } catch (Exception e) {
        }
        super.onDestroy();
    }

}