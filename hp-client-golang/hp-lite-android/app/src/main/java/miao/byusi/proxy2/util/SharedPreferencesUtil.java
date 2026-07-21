package miao.byusi.proxy2.util;

import android.content.Context;
import android.content.SharedPreferences;

public class SharedPreferencesUtil {

    private static final String PREFERENCE_NAME = "config";
    public static final String APK_VERSION = "APK_VERSION";
    public static final String APK_DOWNLOAD_URL = "APK_DOWNLOAD_URL";

    private static final Object lock = new Object();
    private static SharedPreferences sharedPreferences;

    private static SharedPreferences getPreferences(Context context) {
        synchronized (lock) {
            if (sharedPreferences == null) {
                sharedPreferences = context.getSharedPreferences(PREFERENCE_NAME, Context.MODE_PRIVATE);
            }
            return sharedPreferences;
        }
    }

    public static void putBoolean(Context context, String key, boolean value) {
        getPreferences(context).edit().putBoolean(key, value).apply();
    }

    public static boolean getBoolean(Context context, String key, boolean value) {
        return getPreferences(context).getBoolean(key, value);
    }

    public static void putString(Context context, String key, String value) {
        getPreferences(context).edit().putString(key, value).apply();
    }

    public static String getString(Context context, String key, String defValue) {
        return getPreferences(context).getString(key, defValue);
    }

    public static void putInt(Context context, String key, int value) {
        getPreferences(context).edit().putInt(key, value).apply();
    }

    public static int getInt(Context context, String key, int defValue) {
        return getPreferences(context).getInt(key, defValue);
    }

    public static void remove(Context context, String key) {
        getPreferences(context).edit().remove(key).apply();
    }
}