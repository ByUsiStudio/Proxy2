import {createApp} from 'vue'
import App from './App.vue'
import {router} from "./router/index.js";
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

const elementApp = createApp(App);
elementApp.use(Antd)
elementApp.use(router)
elementApp.mount('#app')

// 51.la 统计组件样式
const style = document.createElement('style');
style.textContent = `
  .la-data-widget__container{
    display: block;
    text-align: center;
    background-color: #4b6ff6;
  }
  .la-data-widget__container span{
    color: #ffffff !important;
  }
`;
document.head.appendChild(style);
