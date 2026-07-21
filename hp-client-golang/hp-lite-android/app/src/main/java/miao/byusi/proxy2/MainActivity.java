package miao.byusi.proxy2;

import android.Manifest;
import android.content.DialogInterface;
import android.content.Intent;
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
    private int REQUEST_CODE_SCAN = 111;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        //设备ID
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

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();
        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            show();
            return true;
        }
        if (id == R.id.action_about) {
            Toast.makeText(getApplicationContext(), "禁止一切违法行为，后果自负，开源地址：https://gitee.com/byusistudio/proxy2", Toast.LENGTH_LONG).show();
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
// 设置自定义布局
        LayoutInflater inflater = getLayoutInflater();
        View dialogView = inflater.inflate(R.layout.dialog_layout, null);
        builder.setView(dialogView);

// 获取输入框的引用
        EditText connect_edit = dialogView.findViewById(R.id.connect_edittext);
        //回线
        String connect = SharedPreferencesUtil.getString(getApplicationContext(), ConstConfig.CONNECT, "");
        connect_edit.setText(connect.trim());
// 设置确定按钮
        builder.setPositiveButton("确定", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                String connect = connect_edit.getText().toString().trim();
                if (connect.trim().isEmpty()) {
                    Toast.makeText(getApplicationContext(), "请输入连接码", Toast.LENGTH_LONG).show();
                    return;
                }
                SharedPreferencesUtil.putString(getApplicationContext(), ConstConfig.CONNECT, connect);
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
                            public void onPermissionsChecked(MultiplePermissionsReport multiplePermissionsReport) {
                                Intent intent = new Intent(MainActivity.this, CaptureActivity.class);
                                /*ZxingConfig是配置类
                                 *可以设置是否显示底部布局，闪光灯，相册，
                                 * 是否播放提示音  震动
                                 * 设置扫描框颜色等
                                 * 也可以不传这个参数
                                 * */
                                ZxingConfig config = new ZxingConfig();
                                // config.setPlayBeep(false);//是否播放扫描声音 默认为true
                                //  config.setShake(false);//是否震动  默认为true
                                // config.setDecodeBarCode(false);//是否扫描条形码 默认为true
//                                config.setReactColor(R.color.colorAccent);//设置扫描框四个角的颜色 默认为白色
//                                config.setFrameLineColor(R.color.colorAccent);//设置扫描框边框颜色 默认无色
//                                config.setScanLineColor(R.color.colorAccent);//设置扫描线的颜色 默认白色
                                config.setFullScreenScan(false);//是否全屏扫描  默认为true  设为false则只会在扫描框中扫描
                                intent.putExtra(Constant.INTENT_ZXING_CONFIG, config);
                                startActivityForResult(intent, REQUEST_CODE_SCAN);
                            }
                            @Override
                            public void onPermissionRationaleShouldBeShown(List<PermissionRequest> list, PermissionToken permissionToken) {
                                Toast.makeText(getApplicationContext(), "未获取权限", Toast.LENGTH_LONG).show();
                            }
                        }).check();

            }
        });

// 设置取消按钮
        builder.setNegativeButton("取消", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                dialog.dismiss();
            }
        });
// 创建并显示对话框
        AlertDialog dialog = builder.create();
        dialog.show();
    }


    public void reStartConnect() {
        Intent restartIntent = new Intent(getApplicationContext(), ProxyService.class);
        getApplicationContext().stopService(restartIntent);
        getApplicationContext().startService(restartIntent);

    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);

        // 扫描二维码/条码回传
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