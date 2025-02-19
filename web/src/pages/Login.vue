<template>
    <div class="login-container">
        <div class="login-box">
            <div class="login-header">
                <h2>{{ t('login.welcome') }}</h2>
                <p>{{ t('login.pleaseLogin') }}</p>
            </div>
            <form @submit.prevent="handleLogin">
                <div class="form-group">
                    <div class="input-group">
                        <i class="fas fa-user"></i>
                        <input v-model="form.email" type="text" required placeholder="请输入邮箱">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group">
                        <i class="fas fa-lock"></i>
                        <input v-model="form.password" type="password" required placeholder="请输入密码">
                    </div>
                </div>
                <button type="submit" class="login-btn">
                    <span>{{ t('login.login') }}</span>
                    <i class="fas fa-arrow-right"></i>
                </button>
            </form>
            <div class="additional-links">
                <!-- <a href="#forgot-password">忘记密码？</a>
                <a href="#register">注册新账户</a> -->
            </div>
        </div>
    </div>
</template>
<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const form = reactive({
    email: '',
    password: ''
})

const handleLogin = async () => {
    try {
        const { data } = await api.post('login', form)
        localStorage.setItem('token', data.accessToken)
        router.push('/admin/dashboard')
    } catch {
        // 此处不再需要处理错误，已由全局拦截器处理
    }
}
</script>
<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    overflow-y: auto;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
    background: white;
    padding: 40px;
    border-radius: 20px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 440px;
    transition: all 0.3s ease;
    margin: 40px 20px;
}

.login-header {
    text-align: center;
    margin-bottom: 2rem;
}

.login-header h2 {
    color: #2d3748;
    font-size: 1.8rem;
    margin-bottom: 0.5rem;
}

.login-header p {
    color: #718096;
    font-size: 0.9rem;
}

.form-group {
    margin-bottom: 1.8rem;
}

.input-group {
    position: relative;
    display: flex;
    align-items: center;
    background: #f7fafc;
    border-radius: 8px;
    padding: 0 15px;
    transition: all 0.3s ease;
    border: 2px solid #e2e8f0;
}

.input-group:hover {
    border-color: #c3dafe;
}

.input-group:focus-within {
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.input-group i {
    color: #a0aec0;
    margin-right: 10px;
    font-size: 1rem;
}

input {
    width: 100%;
    padding: 12px 0;
    border: none;
    background: transparent;
    font-size: 1rem;
    color: #2d3748;
    outline: none;
}

input::placeholder {
    color: #a0aec0;
    opacity: 1;
}

button {
    width: 100%;
    padding: 15px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    transition: all 0.3s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

button:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

button:active {
    transform: translateY(0);
}

.login-btn span {
    letter-spacing: 0.5px;
}

.login-btn i {
    font-size: 0.9rem;
}

.additional-links {
    margin-top: 2rem;
    display: flex;
    justify-content: space-between;
    font-size: 0.9rem;
}

.additional-links a {
    color: #718096;
    text-decoration: none;
    transition: color 0.3s ease;
}

.additional-links a:hover {
    color: #667eea;
}

@media (max-width: 480px) {
    .login-box {
        padding: 30px 20px;
        margin: 0;
        border-radius: 0;
        max-width: 100%;
        min-height: auto;
        background: rgba(255, 255, 255, 0.95);
        backdrop-filter: blur(10px);
        box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
        display: flex;
        flex-direction: column;
        justify-content: center;
        border: 1px solid rgba(255, 255, 255, 0.2);
    }

    .login-header h2 {
        font-size: 1.5rem;
        color: #2d3748;
    }

    button {
        padding: 12px;
    }

    .login-container {
        padding: 20px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }
}
</style>