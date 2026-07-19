<template>
  <div>
    <a-tabs default-active-key="basic">
      <a-tab-pane key="basic" tab="基础设置">
        <a-form :model="configForm" layout="vertical">
          <a-form-item label="站点标题">
            <a-input v-model:value="configForm.siteTitle" placeholder="请输入站点标题" />
          </a-form-item>
          <a-form-item label="公开注册">
            <a-switch v-model:checked="configForm.openRegister" />
          </a-form-item>
          <a-form-item label="注册审核">
            <a-switch v-model:checked="configForm.registerReview" />
            <span style="margin-left: 10px; color: #888;">开启后新用户注册需要管理员审核</span>
          </a-form-item>
          <a-form-item>
            <a-button class="btn edit" @click="saveBasic">保存设置</a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="smtp" tab="SMTP邮件配置">
        <a-form :model="smtpForm" layout="vertical">
          <a-form-item label="启用SMTP">
            <a-switch v-model:checked="smtpForm.enabled" />
          </a-form-item>
          <a-form-item label="SMTP服务器地址">
            <a-input v-model:value="smtpForm.host" placeholder="如: smtp.example.com" />
          </a-form-item>
          <a-form-item label="SMTP端口">
            <a-input-number v-model:value="smtpForm.port" :min="1" :max="65535" />
          </a-form-item>
          <a-form-item label="SMTP账号">
            <a-input v-model:value="smtpForm.username" placeholder="请输入SMTP账号" />
          </a-form-item>
          <a-form-item label="SMTP密码">
            <a-input-password v-model:value="smtpForm.password" placeholder="请输入SMTP密码" />
          </a-form-item>
          <a-form-item label="发件人邮箱">
            <a-input v-model:value="smtpForm.from" placeholder="请输入发件人邮箱" />
          </a-form-item>
          <a-form-item label="发件人名称">
            <a-input v-model:value="smtpForm.fromName" placeholder="请输入发件人名称" />
          </a-form-item>
          <a-form-item label="启用SSL">
            <a-switch v-model:checked="smtpForm.enableSSL" />
          </a-form-item>
          <a-form-item>
            <a-button class="btn edit" @click="saveSmtp">保存配置</a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup>
import { onMounted, reactive, watch } from "vue";
import { notification } from "ant-design-vue";
import { getSystemConfig, updateSystemConfig } from "../../api/client/user";

const configForm = reactive({
  siteTitle: '',
  openRegister: false,
  registerReview: false,
});

const smtpForm = reactive({
  enabled: false,
  host: '',
  port: 587,
  username: '',
  password: '',
  from: '',
  fromName: '',
  enableSSL: false,
});

const loadConfig = () => {
  getSystemConfig().then(res => {
    if (res.code === 200 && res.data) {
      const data = res.data;
      configForm.siteTitle = data.siteTitle || '';
      configForm.openRegister = data.openRegister || false;
      configForm.registerReview = data.registerReview || false;
      
      if (data.smtp) {
        smtpForm.enabled = data.smtp.enabled || false;
        smtpForm.host = data.smtp.host || '';
        smtpForm.port = data.smtp.port || 587;
        smtpForm.username = data.smtp.username || '';
        smtpForm.password = data.smtp.password || '';
        smtpForm.from = data.smtp.from || '';
        smtpForm.fromName = data.smtp.fromName || '';
        smtpForm.enableSSL = data.smtp.enableSSL || false;
      }
    }
  });
};

const saveBasic = () => {
  updateSystemConfig({
    siteTitle: configForm.siteTitle,
    openRegister: configForm.openRegister,
    registerReview: configForm.registerReview,
    smtp: {
      enabled: smtpForm.enabled,
      host: smtpForm.host,
      port: smtpForm.port,
      username: smtpForm.username,
      password: smtpForm.password,
      from: smtpForm.from,
      fromName: smtpForm.fromName,
      enableSSL: smtpForm.enableSSL,
    },
  }).then(res => {
    if (res.code === 200) {
      notification.success({
        message: '保存成功',
      });
    }
  });
};

const saveSmtp = () => {
  updateSystemConfig({
    siteTitle: configForm.siteTitle,
    openRegister: configForm.openRegister,
    registerReview: configForm.registerReview,
    smtp: {
      enabled: smtpForm.enabled,
      host: smtpForm.host,
      port: smtpForm.port,
      username: smtpForm.username,
      password: smtpForm.password,
      from: smtpForm.from,
      fromName: smtpForm.fromName,
      enableSSL: smtpForm.enableSSL,
    },
  }).then(res => {
    if (res.code === 200) {
      notification.success({
        message: '保存成功',
      });
    }
  });
};

onMounted(() => {
  loadConfig();
});
</script>