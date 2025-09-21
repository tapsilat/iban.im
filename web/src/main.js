import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import axios from 'axios';
import './tailwind.css';

axios.interceptors.response.use(
    response => response,
    error => {
        if (error?.response?.status === 401) {
            localStorage.removeItem('user');
            location.href = '/login';
        }
        return Promise.reject(error);
    }
);

const app = createApp(App);
app.use(store);
app.use(router);
app.mount('#app');
