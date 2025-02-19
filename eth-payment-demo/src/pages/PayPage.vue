<template>
    <div class="skin-container">
        <div class="payment-container">
            <!-- 皮肤切换器 -->
            <div class="skin-switcher">
                <el-dropdown @command="changeLanguage">
                    <span class="lang-selector"> {{ t('common.language') }}: {{ currentLang }} <i
                            class="el-icon-arrow-down"></i>
                    </span>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item v-for="lang in availableLangs" :key="lang.value" :command="lang.value"> {{
                                lang.label }} </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
                <el-dropdown @command="changeSkin">
                    <i class="icon-skin"></i>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item v-for="(skin, key) in skins" :key="key" :command="key"
                                :class="['skin-item', { active: currentSkin === key }]">
                                <span class="color-preview"
                                    :style="{ background: skin.vars['--primary-color'] }"></span> {{ skin.name }}
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
            <!-- 主卡片容器 -->
            <div class="payment-card" :class="`skin-${currentSkin}`">
                <!-- 订单不存在状态 -->
                <div v-if="orderNotFound" class="state-card error-state">
                    <i class="icon-error"></i>
                    <h2>{{ t('error.orderNotFound') }}</h2>
                    <p>{{ t('error.orderExpired') }}</p>
                    <el-button type="primary" class="full-width-btn" @click="toCheemsHappy"> {{ t('common.backHome') }}
                    </el-button>
                </div>
                <!-- 加载状态 -->
                <div v-else-if="loading" class="state-card loading-state">
                    <div class="skeleton-loader">
                        <div class="skeleton-line long"></div>
                        <div class="skeleton-line"></div>
                        <div class="skeleton-line medium"></div>
                    </div>
                </div>
                <!-- 正常状态 -->
                <template v-else>
                    <!-- 订单信息 -->
                    <div class="order-info">
                        <h2 class="order-title">
                            <i class="icon-order"></i> {{ t('order.paymentOrder') }} #{{ orderInfo.order_no }}
                        </h2>
                        <div class="info-grid">
                            <div class="info-item">
                                <span class="info-label">{{ t('order.merchant') }}</span>
                                <span class="info-value highlight">{{ orderInfo.mer_name }}</span>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('order.amount') }}</span>
                                <span class="info-value highlight">{{ orderInfo.amount }} {{ orderInfo.coin }}</span>
                            </div>
                            <!-- <div class="info-item">
                            <span class="info-label">收款地址</span>
                            <span class="info-value mono">{{ orderInfo.ac }}</span>
                        </div> -->
                            <div class="info-item">
                                <span class="info-label">{{ t('order.targetNetwork') }}</span>
                                <span class="info-value network-tag">{{ chainName }}</span>
                            </div>
                        </div>
                    </div>
                    <!-- 钱包连接区域 -->
                    <div class="wallet-section" :class="{ connected: isConnected }" v-if="!txHash">
                        <template v-if="!isConnected">
                            <el-button type="primary" class="full-width-btn" :loading="connecting"
                                @click="connectWallet">
                                <i class="icon-wallet"></i> {{ t('wallet.connectMetaMask') }} </el-button>
                            <p class="wallet-tip">{{ t('wallet.ensureNetwork', { network: chainName }) }}</p>
                        </template>
                        <div v-else class="wallet-connected">
                            <i class="icon-success"></i>
                            <div class="wallet-info">
                                <p>已连接钱包</p>
                                <p class="wallet-address">{{ shortAddress }}</p>
                            </div>
                        </div>
                    </div>
                    <!-- 支付操作 -->
                    <div v-if="isConnected && !txHash" class="payment-action">
                        <el-button type="success" class="full-width-btn" :loading="processing" @click="handlePayment">
                            <i class="icon-pay"></i> {{ t('common.confirmPayment', {
                                amount: orderInfo.amount, coin:
                                    orderInfo.coin
                            }) }} </el-button>
                    </div>
                    <!-- 交易状态 -->
                    <div v-if="txHash" class="transaction-status">
                        <i class="icon-success"></i>
                        <p class="success-title">{{ t('order.paymentSuccess') }}</p>
                        <div class="countdown-box">
                            <p>{{ t('order.autoRedirect', { seconds: countdown }) }}</p>
                            <el-button type="primary" size="small" @click="toCheemsHappy"> 立即跳转 </el-button>
                        </div>
                    </div>
                </template>
                <!-- 全局错误提示 -->
                <div v-if="errorMessage" class="global-error">
                    <el-alert :title="errorMessage" type="error" show-icon :closable="false" />
                    <el-button v-if="showRetry" class="retry-btn" size="small" @click="handleRetry"> 重试 </el-button>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { createI18n, useI18n } from 'vue-i18n'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ethers } from 'ethers'
import api from '../api'
import { ElMessage } from 'element-plus'
import payZh from '../i18n/pay'
import payEn from '../i18n/pay-en'
import elementZh from 'element-plus/es/locale/lang/zh-cn'
import elementEn from 'element-plus/es/locale/lang/en'

// 创建独立i18n实例
const i18n = createI18n({
    legacy: false,
    locale: localStorage.getItem('payLang') || 'zh',
    messages: {
        zh: { ...payZh, ...elementZh },
        en: { ...payEn, ...elementEn }
    }
})

const { t, locale } = i18n.global
const route = useRoute()

// 订单状态
const orderInfo = ref(null)
const loading = ref(true)
const orderNotFound = ref(false)
const errorMessage = ref('')
const showRetry = ref(false)

// 钱包状态
const isConnected = ref(false)
const connecting = ref(false)
const currentAccount = ref('')
let provider = null

// 支付状态
const processing = ref(false)
const txHash = ref('')

// 支付完成状态和倒计时功能
const countdown = ref(10)
const timer = ref(null)

// 新增皮肤相关代码
const currentSkin = ref('default')
const skins = computed(() => ({
    default: {
        name: t('skin.default'),
        vars: {
            '--primary-color': '#409EFF',
            '--success-color': '#67C23A',
            '--warning-color': '#E6A23C',
            '--danger-color': '#F56C6C',
            '--bg-color': '#ffffff',
            '--card-bg': '#f8fafb',
            '--text-primary': '#303133',
            '--text-regular': '#606266',
            '--border-color': '#e4e7ed',
            '--gradient': 'linear-gradient(120deg, #fdfbfb 0%, #ebedee 100%)',
            '--button-text': '#ffffff',
            '--button-hover': '#337ecc',
            '--input-bg': '#ffffff',
            '--input-text': '#606266'
        }
    },
    dark: {
        name: t('skin.dark'),
        vars: {
            '--primary-color': '#409EFF',
            '--bg-color': '#121212',
            '--card-bg': '#2d2d2d',
            '--text-primary': '#ffffff',
            '--text-regular': '#e0e0e0',
            '--border-color': '#434343',
            '--gradient': 'linear-gradient(145deg, #2c3e50 0%, #1a1a1a 100%)',
            '--button-text': '#121212',
            '--button-hover': '#66b1ff',
            '--input-bg': '#2d2d2d',
            '--input-text': '#e0e0e0',
            '--link-color': '#66b1ff'
        }
    },
    nature: {
        name: 'Nature',
        vars: {
            '--primary-color': '#48c78e',
            '--success-color': '#3ec46d',
            '--warning-color': '#ffe08a',
            '--danger-color': '#f14668',
            '--bg-color': '#f5f7fa',
            '--card-bg': '#ffffff',
            '--text-primary': '#363636',
            '--text-regular': '#4a4a4a',
            '--border-color': '#dbdbdb',
            '--gradient': 'linear-gradient(120deg, #d4fc79 0%, #96e6a1 100%)',
            '--button-text': '#ffffff',
            '--button-hover': '#5daf34',
            '--input-bg': '#f5f7fa',
            '--input-text': '#363636'
        }
    },
    purple: {
        name: 'Purple',
        vars: {
            '--primary-color': '#9c27b0',
            '--success-color': '#4caf50',
            '--warning-color': '#ff9800',
            '--danger-color': '#f44336',
            '--bg-color': '#f3e5f5',
            '--card-bg': '#ffffff',
            '--text-primary': '#4a148c',
            '--text-regular': '#6a1b9a',
            '--border-color': '#d1c4e9',
            '--gradient': 'linear-gradient(120deg, #e1bee7 0%, #ce93d8 100%)',
            '--button-text': '#ffffff',
            '--button-hover': '#6a1b9a',
            '--input-bg': '#f3e5f5',
            '--input-text': '#4a148c'
        }
    },
    orange: {
        name: 'Orange',
        vars: {
            '--primary-color': '#ff9800',
            '--success-color': '#8bc34a',
            '--warning-color': '#ffc107',
            '--danger-color': '#f44336',
            '--bg-color': '#fff3e0',
            '--card-bg': '#ffffff',
            '--text-primary': '#ef6c00',
            '--text-regular': '#fb8c00',
            '--border-color': '#ffcc80',
            '--gradient': 'linear-gradient(120deg, #ffe0b2 0%, #ffcc80 100%)'
        }
    },
    pink: {
        name: 'Pink',
        vars: {
            '--primary-color': '#e91e63',
            '--success-color': '#4caf50',
            '--warning-color': '#ffc107',
            '--danger-color': '#f44336',
            '--bg-color': '#fce4ec',
            '--card-bg': '#ffffff',
            '--text-primary': '#ad1457',
            '--text-regular': '#c2185b',
            '--border-color': '#f8bbd0',
            '--gradient': 'linear-gradient(120deg, #f8bbd0 0%, #f48fb1 100%)'
        }
    }
}))

const toCheemsHappy = () => {
    if (orderInfo && orderInfo.value?.return_url) {
        window.location.href = orderInfo.value.return_url
    } else {
        window.location.href = 'https://cheemshappy.com'
    }
}

// 获取订单信息
const fetchOrder = async () => {
    try {
        loading.value = true
        errorMessage.value = ''
        const { data } = await api.get(`/order/order/${route.query.no}`)
        orderInfo.value = data
    } catch (err) {
        handleOrderError(err)
    } finally {
        loading.value = false
        showRetry.value = !!errorMessage.value
    }
}

// 错误处理
const handleOrderError = (err) => {
    orderInfo.value = null
    if (err.code === 404) {
        orderNotFound.value = true
    } else {
        errorMessage.value = err.response?.data?.message || '网络连接异常，请稍后重试'
    }
}

// 钱包连接
const connectWallet = async () => {
    try {
        connecting.value = true
        errorMessage.value = ''

        if (!window.ethereum) {
            throw new Error('请安装MetaMask钱包')
        }

        provider = new ethers.BrowserProvider(window.ethereum)
        const accounts = await provider.send("eth_requestAccounts", [])

        if (accounts.length > 0) {
            currentAccount.value = accounts[0]
            isConnected.value = true
            setupNetworkListener()
        }
    } catch (err) {
        errorMessage.value = err.code === 4001
            ? '用户取消钱包连接'
            : `钱包连接失败：${err.message || '未知错误'}`
        showRetry.value = true
    } finally {
        connecting.value = false
    }
}

// 支付处理
const handlePayment = async () => {
    try {
        processing.value = true
        errorMessage.value = ''

        const signer = await provider.getSigner()
        const network = await provider.getNetwork()

        // 网络检查
        if (network.chainId !== BigInt(orderInfo.value.chain)) {
            await switchNetwork()
        }

        // 执行支付
        const tx = await signer.sendTransaction({
            to: orderInfo.value.ac,
            value: ethers.parseEther(orderInfo.value.amount),
            chainId: parseInt(orderInfo.value.chain)
        })

        txHash.value = tx.hash
        await api.post(`/order/${orderInfo.value.order_no}/tx`, { txHash: tx.hash })

        ElMessage.success('支付成功')
        startCountdown() // 启动倒计时

    } catch (err) {
        handlePaymentError(err)
    } finally {
        processing.value = false
    }
}

// 网络切换
const switchNetwork = async () => {
    try {
        await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: `0x${parseInt(orderInfo.value.chain).toString(16)}` }]
        })
    } catch (err) {
        if (err.code === 4001) throw new Error('用户取消网络切换')
        throw err
    }
}

// 支付错误处理
const handlePaymentError = (err) => {
    errorMessage.value = err.code === 4001
        ? '用户取消交易'
        : `支付失败：${err.message || '区块链网络异常'}`
    showRetry.value = true
}

// 重试操作
const handleRetry = () => {
    errorMessage.value = ''
    showRetry.value = false
    if (orderInfo.value) {
        handlePayment()
    } else {
        fetchOrder()
    }
}

// 计算属性
const shortAddress = computed(() => {
    const addr = currentAccount.value
    return addr ? `${addr.slice(0, 6)}...${addr.slice(-4)}` : ''
})

const chainName = computed(() => {
    const chains = {
        1: 'Ethereum',
        56: 'BNB Chain',
        137: 'Polygon',
        43114: 'Avalanche',
        10: 'Optimism',
        42161: 'Arbitrum',
        5: 'Goerli Testnet',
        80001: 'Mumbai Testnet',
        11155111: 'Sepolia Testnet',
    }
    return chains[orderInfo.value?.chain] || `Chain ID: ${orderInfo.value?.chain}`
})

// 监听网络变化
const setupNetworkListener = () => {
    if (window.ethereum) {
        window.ethereum.on('chainChanged', (chainId) => {
            if (parseInt(chainId) !== parseInt(orderInfo.value?.chain)) {
                ElMessage.warning('检测到网络切换，请切换回目标网络')
                isConnected.value = false
            }
        })
    }
}

// 在组件卸载时清理
onUnmounted(() => {
    if (window.ethereum) {
        window.ethereum.removeAllListeners('chainChanged')
    }
    if (timer.value) clearInterval(timer.value)
})

// 在支付成功时启动倒计时
const startCountdown = () => {
    countdown.value = 10
    timer.value = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
            if (orderInfo.value.return_url) {
                console.log('跳转', orderInfo.value.return_url)
                window.location.href = orderInfo.value.return_url
            } else {
                window.location.href = 'https://cheemshappy.com'
            }
        }
    }, 1000)
}

// 修改皮肤应用逻辑
const applySkin = (skinKey) => {
    const skinConfig = skins.value[skinKey]
    if (!skinConfig || !skinConfig.vars) {
        return
    }

    // 应用CSS变量
    const root = document.documentElement
    Object.entries(skinConfig.vars).forEach(([key, value]) => {
        root.style.setProperty(key, value)
    })
}

// 修改皮肤切换逻辑
const changeSkin = (skinKey) => {
    currentSkin.value = skinKey
    localStorage.setItem('selectedSkin', skinKey)
    applySkin(skinKey)
}

// 确保初始化时使用有效皮肤
onMounted(() => {
    let savedSkin = localStorage.getItem('selectedSkin')
    if (!savedSkin || !skins.value[savedSkin]) {
        const skinKeys = Object.keys(skins.value)
        savedSkin = skinKeys[Math.floor(Math.random() * skinKeys.length)]
    }
    changeSkin(savedSkin)
    fetchOrder()
})

const currentLang = computed(() => locale.value === 'zh' ? '中文' : 'EN')

const changeLanguage = (lang) => {
    i18n.global.locale.value = lang
    localStorage.setItem('payLang', lang)
}

const availableLangs = [
    { value: 'zh', label: t('common.zh') },
    { value: 'en', label: t('common.en') }
]
</script>
<style scoped>
.skin-container {
    display: flex;
    width: 100%;
    background-color: var(--bg-color);
}

.payment-container {
    max-width: 480px;
    margin: 2rem auto;
    padding: 0 1rem;
    background: var(--bg-color);
    min-height: 100vh;
    padding: 20px;
    transition: background 0.3s ease;
}

.payment-card {
    background: var(--card-bg);
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    padding: 2rem;
    position: relative;
    border: 1px solid var(--border-color);
}

.state-card {
    text-align: center;
    padding: 2rem 0;

    h2 {
        color: #1a1a1a;
        margin: 1.5rem 0 0.5rem;
        font-size: 1.25rem;
    }

    p {
        color: #666;
        margin-bottom: 1.5rem;
    }
}

.order-info {
    margin-bottom: 2rem;

    .order-title {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: #1a1a1a;
        margin-bottom: 1.5rem;
        font-size: 1.25rem;
    }
}

.info-grid {
    display: grid;
    gap: 1rem;

    .info-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        background: var(--card-bg);
        border-radius: 8px;
        border: 1px solid var(--border-color);

        .info-label {
            color: #666;
            font-size: 0.8rem;
            opacity: 0.8;
            transition: opacity 0.3s ease;
        }

        .info-item:hover .info-label {
            opacity: 1;
        }

        .info-value {
            color: #1a1a1a;
            font-weight: 500;

            &.highlight {
                color: var(--primary-color);
                font-weight: 600;
            }

            &.mono {
                font-family: monospace;
            }
        }

        .network-tag {
            background: #e6f7ff;
            color: #1890ff;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.85rem;
        }
    }
}

.wallet-section {
    margin: 2rem 0;

    &.connected {
        background: var(--card-bg);
        border-color: var(--border-color);
        border-radius: 8px;
        padding: 1rem;
        transition: all 0.3s ease;
    }

    .wallet-connected {
        display: flex;
        align-items: center;
        gap: 1rem;

        i {
            color: #3b82f6;
            font-size: 1.5rem;
        }

        .wallet-info {
            flex-grow: 1;

            p:first-child {
                color: #1e293b;
                font-weight: 500;
            }

            .wallet-address {
                color: #64748b;
                font-size: 0.9rem;
            }
        }
    }

    .wallet-tip {
        color: #666;
        font-size: 0.85rem;
        margin-top: 0.5rem;
        text-align: center;
    }
}

.full-width-btn {
    width: 100%;
    padding: 1rem;
    font-size: 1rem;
}

.transaction-status {
    text-align: center;
    margin-top: 1.5rem;
    padding: 1rem;
    background: var(--card-bg);
    border-radius: 8px;
    border: 1px solid var(--success-color);

    .success-title {
        color: #52c41a;
        font-weight: 500;
        margin: 0.5rem 0;
    }

    .countdown-box {
        margin-top: 1rem;

        p {
            color: #666;
            margin-bottom: 0.5rem;
        }

        .el-button {
            padding: 0.5rem 1.5rem;
        }
    }
}

.global-error {
    margin-top: 1.5rem;

    .retry-btn {
        margin-top: 1rem;
        width: 100%;
    }
}

/* 图标样式 */
[class^="icon-"] {
    width: 24px;
    height: 24px;
    display: inline-block;
    background-size: contain;
}

.icon-error {
    background-image: url('data:image/svg+xml;utf8,<svg ...>/* 错误图标SVG代码 */</svg>');
}

.icon-success {
    background-image: url('data:image/svg+xml;utf8,<svg ...>/* 成功图标SVG代码 */</svg>');
}

/* 骨架屏加载动画 */
@keyframes shimmer {
    0% {
        background-position: -200% 0;
    }

    100% {
        background-position: 200% 0;
    }
}

.skeleton-loader {
    padding: 1rem;

    .skeleton-line {
        height: 20px;
        background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
        background-size: 200% 100%;
        border-radius: 4px;
        margin-bottom: 1rem;
        animation: shimmer 1.5s infinite;

        &.long {
            width: 80%;
        }

        &.medium {
            width: 60%;
        }
    }
}

/* 新增皮肤切换器样式 */
.skin-switcher {
    display: flex;
    justify-content: center;
    align-items: center;
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 1000;

    .icon-skin {
        display: block;
        width: 32px;
        height: 32px;
        background: var(--primary-color);
        border-radius: 50%;
        cursor: pointer;
        transition: transform 0.3s;

        &:hover {
            transform: rotate(180deg);
        }
    }
}

/* 不同皮肤的特效 */
.skin-dark {
    .payment-card {
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .info-item {
        background: rgba(255, 255, 255, 0.08);
    }

    .transaction-status {
        background: rgba(40, 40, 40, 0.9);
        border-color: var(--success-color);
    }

    .info-label {
        color: rgba(255, 255, 255, 0.7);
    }

    .el-alert--error {
        background: rgba(245, 108, 108, 0.1) !important;
    }

    .skeleton-loader .skeleton-line {
        background: linear-gradient(90deg, #2d2d2d 25%, #404040 50%, #2d2d2d 75%);
    }
}

.skin-nature {
    .payment-card {
        background-image: var(--gradient);
        border: 1px solid var(--success-color);
    }

    .el-button--primary {
        background: var(--primary-color);
        border-color: var(--primary-color);

        &:hover {
            opacity: 0.9;
        }
    }

    .el-button--success {
        background: var(--success-color) !important;
    }
}

/* 下拉菜单样式 */
.skin-item {
    display: flex;
    align-items: center;
    padding: 8px 12px;

    &.active {
        background: rgba(var(--primary-color), 0.1);
    }

    .color-preview {
        width: 16px;
        height: 16px;
        border-radius: 4px;
        margin-right: 8px;
        border: 1px solid var(--border-color);
    }
}

/* 所有按钮增加过渡效果 */
.el-button {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
    color: var(--button-text) !important;
    background: var(--primary-color) !important;
    border-color: var(--primary-color) !important;

    &:hover {
        background: var(--button-hover) !important;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
}

/* 输入框增加主题色边框 */
.el-input__wrapper {
    border-color: var(--border-color) !important;
    background: var(--input-bg) !important;
    color: var(--input-text) !important;
}

/* 链接颜色 */
a {
    color: var(--link-color, var(--primary-color));

    &:hover {
        opacity: 0.8;
    }
}

.lang-selector {
    margin-right: 15px;
    cursor: pointer;
    color: var(--text-primary);
    font-weight: 500;
}
</style>