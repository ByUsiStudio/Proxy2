package miao.byusi.proxy2;

import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.os.Build;
import android.os.Handler;
import android.os.IBinder;
import android.os.Looper;
import android.os.Message;
import android.os.PowerManager;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.core.app.NotificationCompat;

import java.lang.ref.WeakReference;
import java.text.SimpleDateFormat;
import java.util.Date;

import hp_android_lib.Callback;
import hp_android_lib.Hp_android_lib;
import miao.byusi.proxy2.util.ConstConfig;
import miao.byusi.proxy2.util.SharedPreferencesUtil;

public class ProxyService extends Service {
    private static final int NOTIFICATION_ID = 1;
    private static final String channelId = "Proxy2Channel";
    private static final String channelName = "Proxy2 Service";
    private NotificationManager notificationManager;
    private static boolean isStart = false;

    private StatusHandler handler;
    private Runnable runnableCode;
    private PowerManager.WakeLock wakeLock;

    @Override
    public void onCreate() {
        Log.i("Proxy2", "Service onCreate - Thread ID = " + Thread.currentThread().getId());
        super.onCreate();

        handler = new StatusHandler(this);

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            int importance = NotificationManager.IMPORTANCE_DEFAULT;
            NotificationChannel channel = new NotificationChannel(channelId, channelName, importance);
            notificationManager = getSystemService(NotificationManager.class);
            notificationManager.createNotificationChannel(channel);
        } else {
            notificationManager = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
        }

        Intent notificationIntent = new Intent(this, MainActivity.class);
        PendingIntent pendingIntent = PendingIntent.getActivity(this, 0, notificationIntent,
                Build.VERSION.SDK_INT >= Build.VERSION_CODES.M ? PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE : PendingIntent.FLAG_UPDATE_CURRENT);

        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, channelId)
                .setSmallIcon(R.drawable.ic_launcher)
                .setContentTitle("Proxy2")
                .setContentText("服务已经运行")
                .setContentIntent(pendingIntent)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);
        startForeground(NOTIFICATION_ID, builder.build());

        PowerManager powerManager = (PowerManager) getSystemService(Context.POWER_SERVICE);
        wakeLock = powerManager.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "Proxy2:WakeLock");
        wakeLock.acquire();

        runnableCode = new Runnable() {
            @Override
            public void run() {
                new Thread(() -> {
                    try {
                        boolean status = Hp_android_lib.getStatus();
                        handler.sendMessage(getMessage(String.valueOf(status)));
                    } catch (Exception e) {
                        Log.e("Proxy2", "getStatus error: ", e);
                        handler.sendMessage(getMessage("false"));
                    }
                }).start();
                handler.postDelayed(this, 1000 * 10);
            }
        };
    }

    private void updateNotification(String message) {
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, channelId)
                .setContentTitle("Proxy2")
                .setContentText(message)
                .setSmallIcon(R.drawable.ic_launcher)
                .setPriority(NotificationCompat.PRIORITY_LOW);
        notificationManager.notify(NOTIFICATION_ID, builder.build());
    }

    public Message getMessage(String msg) {
        Message message = new Message();
        message.what = -1;
        message.obj = msg;
        return message;
    }

    public static class StatusHandler extends Handler {
        private final WeakReference<ProxyService> serviceRef;

        public StatusHandler(ProxyService service) {
            super(Looper.getMainLooper());
            this.serviceRef = new WeakReference<>(service);
        }

        private String getTime() {
            Date currentDate = new Date();
            SimpleDateFormat sdf = new SimpleDateFormat("HH:mm:ss");
            return sdf.format(currentDate);
        }

        @Override
        public void handleMessage(@NonNull Message msg) {
            ProxyService service = serviceRef.get();
            if (service != null && msg.obj != null) {
                String status = msg.obj.toString();
                if ("true".equals(status)) {
                    service.updateNotification("服务连接正常-检查时间:" + getTime());
                } else {
                    service.updateNotification("服务连接异常-检查时间:" + getTime());
                }
            }
        }
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        try {
            if (!isStart) {
                isStart = true;
                String connect = SharedPreferencesUtil.getString(getApplicationContext(), ConstConfig.CONNECT, "");
                Log.i("Proxy2", "Start command with connect: " + connect);

                new Thread(() -> {
                    try {
                        Hp_android_lib.start(connect, new Callback() {
                            @Override
                            public void sendResult(String s) {
                                Log.i("Proxy2", "Callback result: " + s);
                                sendBroadcastMessage(s);
                            }
                        });
                    } catch (Exception e) {
                        Log.e("Proxy2", "Hp_android_lib.start error: ", e);
                        sendBroadcastMessage("启动失败: " + e.getMessage());
                    }
                }).start();

                handler.postDelayed(runnableCode, 1000 * 10);
            }
        } catch (Throwable e) {
            Log.e("Proxy2", "onStartCommand error: ", e);
        }
        return START_STICKY;
    }

    private void sendBroadcastMessage(String message) {
        Intent broadcastIntent = new Intent("miao.byusi.proxy2.LOG_MESSAGE");
        broadcastIntent.putExtra("message", message);
        sendBroadcast(broadcastIntent);
    }

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    @Override
    public void onDestroy() {
        isStart = false;
        if (handler != null) {
            handler.removeCallbacksAndMessages(null);
        }
        if (wakeLock != null && wakeLock.isHeld()) {
            wakeLock.release();
        }
        try {
            Hp_android_lib.close();
        } catch (Exception e) {
            Log.e("Proxy2", "Hp_android_lib.close error: ", e);
        }
        super.onDestroy();
    }
}