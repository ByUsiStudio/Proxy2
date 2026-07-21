package miao.byusi.proxy2;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;


import miao.byusi.proxy2.databinding.FragmentFirstBinding;
import miao.byusi.proxy2.util.ConstConfig;
import miao.byusi.proxy2.util.CopyUtil;
import miao.byusi.proxy2.util.SharedPreferencesUtil;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Locale;

public class FirstFragment extends Fragment {
    private static LogAdapter adapter;
    private FragmentFirstBinding binding;

    @Override
    public View onCreateView(
            LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState
    ) {
        binding = FragmentFirstBinding.inflate(inflater, container, false);
        return binding.getRoot();
    }

    public void startConnect() {
        String connect = SharedPreferencesUtil.getString(getActivity(), ConstConfig.CONNECT, "");
        if (connect.trim().isEmpty()) {
            return;
        }
        Intent intent = new Intent(getActivity(), ProxyService.class);
        getActivity().startService(intent);
    }

    public void onViewCreated(@NonNull View view, Bundle savedInstanceState) {
        super.onViewCreated(view, savedInstanceState);
        // 设置适配器
        adapter = new LogAdapter(getActivity(), new ArrayList<>());
        binding.listView.setAdapter(adapter);
        startConnect();
    }

    @Override
    public void onDestroyView() {
        super.onDestroyView();
        binding = null;
    }

    public static class RegHandler extends Handler {
        @Override
        public void handleMessage(@NonNull Message msg) {
            adapter.addLogEntry(new LogEntry(new Date(), msg.obj.toString()));
        }
    }

    // 自定义适配器类
    private class LogAdapter extends ArrayAdapter<LogEntry> {
        private LayoutInflater inflater;
        private List<LogEntry> logEntries;
        private static final int MAX_LOG_COUNT = 100;

        public LogAdapter(Context context, List<LogEntry> logEntries) {
            super(context, 0, logEntries);
            inflater = LayoutInflater.from(context);
            this.logEntries = logEntries;
        }

        @Override
        public View getView(int position, View convertView, ViewGroup parent) {
            if (convertView == null) {
                convertView = inflater.inflate(R.layout.list_item_log, parent, false);
            }
            // 获取当前位置的日志条目
            LogEntry logEntry = getItem(position);

            // 设置时间戳和内容
            TextView textViewTimestamp = convertView.findViewById(R.id.textViewTimestamp);
            TextView textViewContent = convertView.findViewById(R.id.textViewContent);
            textViewContent.setTextSize(10);

            SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault());
            String timestamp = dateFormat.format(logEntry.getTimestamp());

            textViewTimestamp.setText(timestamp);
            textViewContent.setText(logEntry.getContent());


            textViewContent.setText(logEntry.getContent());

            // 为列表项添加点击监听器
            convertView.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    // 在这里执行你想要的操作，比如显示详细信息或处理点击事件
                    LogEntry logEntry1 = logEntries.get(position);
                    if (logEntry1 != null) {
                        CopyUtil.copy(logEntry1.getContent(), getActivity(), true);
                    }
                }
            });
            return convertView;
        }

        public void addLogEntry(LogEntry logEntry) {
            if (logEntries.size() >= MAX_LOG_COUNT) {
                // 如果达到最大条目数，移除最早的一条日志
                logEntries.remove(0);
            }
            logEntries.add(0, logEntry);
            notifyDataSetChanged();
        }
    }

}