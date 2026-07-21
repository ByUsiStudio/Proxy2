package net.hserver.hplite.util;

import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.widget.Toast;

public class CopyUtil {

    public static void copy(String msg, Context context,boolean show) {
        // 获取剪贴板管理器实例
        ClipboardManager clipboard = (ClipboardManager) context.getSystemService(Context.CLIPBOARD_SERVICE);
        // 将文本放入剪贴板
        ClipData clip = ClipData.newPlainText("label", msg);
        clipboard.setPrimaryClip(clip);
        if (show) {
            Toast.makeText(context, "已复制到剪贴板", Toast.LENGTH_SHORT).show();
        }
    }
}
