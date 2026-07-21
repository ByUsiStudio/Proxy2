package miao.byusi.proxy2;

import android.Manifest;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Build;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.Toast;

import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;
import androidx.navigation.NavController;
import androidx.navigation.Navigation;
import androidx.navigation.ui.AppBarConfiguration;
import androidx.navigation.ui.NavigationUI;

import com.google.android.material.snackbar.Snackbar;
import com.karumi.dexter.Dexter;
import com.karumi.dexter.MultiplePermissionsReport;
import com.karumi.dexter.PermissionToken;
import com.karumi.dexter.listener.PermissionRequest;
import com.karumi.dexter.listener.multi.MultiplePermissionsListener;
import com.yzq.zxinglibrary.android.CaptureActivity;
import com.yzq.zxinglibrary.bean.ZxingConfig;
import com.yzq.zxinglibrary.common.Constant;

import miao.byusi.proxy2.databinding.ActivityMainBinding;
import miao.byusi.proxy2.util.ConstConfig;
import miao.byusi.proxy2.util.CopyUtil;
import miao.byusi.proxy2.util.SharedPreferencesUtil;

import java.util.List;

public class MainActivity extends AppCompatActivity {

    private AppBarConfiguration appBarConfiguration;
    private ActivityMainBinding binding;
    private final int REQUEST_CODE_SCAN = 111;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        requestAllPermissions();

        // 设备ID检查
        String connect = SharedPreferencesUtil.getString(getApplicationContext(), ConstConfig.CONNECT, "");
        if (connect == null || connect.trim().isEmpty()) {
            show();
        }

        binding = ActivityMainBinding.inflate(getLayoutInflater());
        setContentView(binding.getRoot());
        setSupportActionBar(binding.toolbar);
        NavController navController = Navigation.findNavController(this, R.id.nav_host_fragment_content_main);
        appBarConfiguration = new AppBarConfiguration.Builder(navController.getGraph()).build();
        NavigationUI.setupActionBarWithNavController(this, navController, appBarConfiguration);

        binding.fab.setOnClickListener(view -> {
            String connect1 = SharedPreferencesUtil.getString(getApplicationContext(), ConstConfig.CONNECT, "");
            CopyUtil.copy(connect1, getApplicationContext(), false);
            Snackbar.make(view, "你的连接码【" + connect1 + "】已复制", Snackbar.LENGTH_LONG)
                    .setAction("Action", null).show();
        });
    }

    private void requestAllPermissions() {
        String[] permissions;
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            // Android 13+ 需要 POST_NOTIFICATIONS 权限
            permissions = new String[]{
                    Manifest.permission.CAMERA,
                    Manifest.permission.READ_EXTERNAL_STORAGE,
                    Manifest.permission.POST_NOTIFICATIONS
            };
        } else {
            permissions = new String[]{
                    Manifest.permission.CAMERA,
                    Manifest.permission.READ_EXTERNAL_STORAGE
            };
        }

        Dexter.withContext(this)
                .withPermissions(permissions)
                .withListener(new MultiplePermissionsListener() {
                    @Override
                    public void onPermissionsChecked(MultiplePermissionsReport report) {
                        if (!report.areAllPermissionsGranted()) {
                            Toast.makeText(MainActivity.this,
                                    "部分权限被拒绝，某些功能可能无法使用",
                                    Toast.LENGTH_LONG).show();
                        }
                    }

                    @Override
                    public void onPermissionRationaleShouldBeShown(List<PermissionRequest> list,
                                                                   PermissionToken permissionToken) {
                        permissionToken.continuePermissionRequest();
                    }
                })
                .onSameThread()
                .check();
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        int id = item.getItemId();
        if (id == R.id.action_settings) {
            show();
            return true;
        }
        if (id == R.id.action_about) {
            Toast.makeText(getApplicationContext(),
                    "禁止一切违法行为，后果自负，开源地址：https://gitee.com/byusistudio/proxy2",
                    Toast.LENGTH_LONG).show();
            return true;
        }
        return super.onOptionsItemSelected(item);
    }

    @Override
    public boolean onSupportNavigateUp() {
        NavController navController = Navigation.findNavController(this, R.id.nav_host_fragment_content_main);
        return NavigationUI.navigateUp(navController, appBarConfiguration)
                || super.onSupportNavigateUp();
    }

    public void show() {
        AlertDialog.Builder builder = new AlertDialog.Builder(this);
        builder.setTitle("设备ID(如果没有请去自建后台添加)");

        LayoutInflater inflater = getLayoutInflater();
        View dialogView = inflater.inflate(R.layout.dialog_layout, null);
        builder.setView(dialogView);

        EditText connect_edit = dialogView.findViewById(R.id.connect_edittext);
        String connect = SharedPreferencesUtil.getString(getApplicationContext(), ConstConfig.CONNECT, "");
        connect_edit.setText(connect.trim());

        // 确定按钮
        builder.setPositiveButton("确定", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                String connect1 = connect_edit.getText().toString().trim();
                if (connect1.trim().isEmpty()) {
                    Toast.makeText(getApplicationContext(), "请输入连接码", Toast.LENGTH_LONG).show();
                    return;
                }
                SharedPreferencesUtil.putString(getApplicationContext(), ConstConfig.CONNECT, connect1);
                reStartConnect();
            }
        });

        builder.setNeutralButton("扫描连接码", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialogInterface, int i) {
                Dexter.withContext(getApplicationContext())
                        .withPermissions(Manifest.permission.CAMERA, Manifest.permission.READ_EXTERNAL_STORAGE)
                        .withListener(new MultiplePermissionsListener() {
                            @Override
                            public void onPermissionsChecked(MultiplePermissionsReport report) {
                                if (report.areAllPermissionsGranted()) {
                                    Intent intent = new Intent(MainActivity.this, CaptureActivity.class);
                                    ZxingConfig config = new ZxingConfig();
                                    config.setFullScreenScan(false);
                                    intent.putExtra(Constant.INTENT_ZXING_CONFIG, config);
                                    startActivityForResult(intent, REQUEST_CODE_SCAN);
                                } else {
                                    Toast.makeText(MainActivity.this,
                                            "需要相机权限才能扫描二维码",
                                            Toast.LENGTH_LONG).show();
                                }
                            }

                            @Override
                            public void onPermissionRationaleShouldBeShown(List<PermissionRequest> list,
                                                                           PermissionToken permissionToken) {
                                permissionToken.continuePermissionRequest();
                                Toast.makeText(getApplicationContext(),
                                        "需要相机权限来扫描二维码",
                                        Toast.LENGTH_LONG).show();
                            }
                        })
                        .onSameThread()
                        .check();
            }
        });

        // 取消按钮
        builder.setNegativeButton("取消", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                dialog.dismiss();
            }
        });

        AlertDialog dialog = builder.create();
        dialog.show();
    }

    public void reStartConnect() {
        Intent restartIntent = new Intent(getApplicationContext(), ProxyService.class);
        // 先停止旧服务
        getApplicationContext().stopService(restartIntent);

        // Android 8.0+ 必须使用 startForegroundService
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            getApplicationContext().startForegroundService(restartIntent);
        } else {
            getApplicationContext().startService(restartIntent);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);

        if (requestCode == REQUEST_CODE_SCAN && resultCode == RESULT_OK) {
            if (data != null) {
                String content = data.getStringExtra(Constant.CODED_CONTENT);
                Toast.makeText(getApplicationContext(), "扫描结果为：" + content, Toast.LENGTH_LONG).show();
                SharedPreferencesUtil.putString(getApplicationContext(), ConstConfig.CONNECT, content);
                show();
            }
        }
    }
}