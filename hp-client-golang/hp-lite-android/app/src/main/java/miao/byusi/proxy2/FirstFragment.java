package miao.byusi.proxy2;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.os.Message;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;

import java.lang.ref.WeakReference;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Locale;

import miao.byusi.proxy2.databinding.FragmentFirstBinding;
import miao.byusi.proxy2.util.ConstConfig;
import miao.byusi.proxy2.util.CopyUtil;
import miao.byusi.proxy2.util.SharedPreferencesUtil;

public class FirstFragment extends Fragment {
    private LogAdapter adapter;
    private FragmentFirstBinding binding;
    private RegHandler regHandler;
    private BroadcastReceiver logReceiver;

    @Override
    public View onCreateView(
            LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState
    ) {
        binding = FragmentFirstBinding.inflate(inflater, container, false);
        regHandler = new RegHandler(this);
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
        adapter = new LogAdapter(getActivity(), new ArrayList<>());
        binding.listView.setAdapter(adapter);
        startConnect();

        logReceiver = new BroadcastReceiver() {
            @Override
            public void onReceive(Context context, Intent intent) {
                String message = intent.getStringExtra("message");
                if (message != null) {
                    addLog(message);
                }
            }
        };
        IntentFilter filter = new IntentFilter("miao.byusi.proxy2.LOG_MESSAGE");
        getActivity().registerReceiver(logReceiver, filter);
    }

    @Override
    public void onDestroyView() {
        super.onDestroyView();
        if (regHandler != null) {
            regHandler.removeCallbacksAndMessages(null);
        }
        if (logReceiver != null) {
            getActivity().unregisterReceiver(logReceiver);
        }
        binding = null;
    }

    public RegHandler getRegHandler() {
        return regHandler;
    }

    public void addLog(String message) {
        if (adapter != null) {
            adapter.addLogEntry(new LogEntry(new Date(), message));
        }
    }

    public static class RegHandler extends Handler {
        private final WeakReference<FirstFragment> fragmentRef;

        public RegHandler(FirstFragment fragment) {
            super(Looper.getMainLooper());
            this.fragmentRef = new WeakReference<>(fragment);
        }

        @Override
        public void handleMessage(@NonNull Message msg) {
            FirstFragment fragment = fragmentRef.get();
            if (fragment != null && fragment.isAdded() && msg.obj != null) {
                fragment.addLog(msg.obj.toString());
            }
        }
    }

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
            LogEntry logEntry = getItem(position);
            if (logEntry == null) {
                return convertView;
            }

            TextView textViewTimestamp = convertView.findViewById(R.id.textViewTimestamp);
            TextView textViewContent = convertView.findViewById(R.id.textViewContent);
            textViewContent.setTextSize(10);

            SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault());
            String timestamp = dateFormat.format(logEntry.getTimestamp());

            textViewTimestamp.setText(timestamp);
            textViewContent.setText(logEntry.getContent());

            convertView.setOnClickListener(v -> {
                LogEntry logEntry1 = logEntries.get(position);
                if (logEntry1 != null && getActivity() != null) {
                    CopyUtil.copy(logEntry1.getContent(), getActivity(), true);
                }
            });
            return convertView;
        }

        public void addLogEntry(LogEntry logEntry) {
            if (logEntries.size() >= MAX_LOG_COUNT) {
                logEntries.remove(0);
            }
            logEntries.add(0, logEntry);
            notifyDataSetChanged();
        }
    }
}