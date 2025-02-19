<template>
    <div class="sys-wallet-settings">
        <el-card>
            <template #header>
                <div class="card-header">
                    <el-button type="primary" size="small" @click="handleCreate" v-if="!walletData">
                        <i class="el-icon-plus"></i> {{ t('sysWallet.create') }} </el-button>
                </div>
            </template>
            <el-alert v-if="walletData" type="error" :closable="false" class="security-alert" show-icon>
                <template #title>
                    <div class="security-title">
                        <div class="security-tips">
                            <p>{{ t('sysWallet.securityWarning') }}</p>
                            <p class="tip-text">（{{ t('sysWallet.longPressTip') }}）</p>
                        </div>
                    </div>
                </template>
            </el-alert>
            <el-empty v-if="!walletData" :description="t('sysWallet.noWallet')" />
            <div v-else class="wallet-info">
                <el-descriptions :column="1" border>
                    <el-descriptions-item :label="t('sysWallet.mnemonic')">
                        <div class="sensitive-field" v-longpress="() => handleLongPress('mnemonic')">
                            <span v-if="showMnemonic" class="mnemonic-text"> {{ walletData.mnemonic }} <el-tag
                                    type="danger" size="mini" class="countdown-tag"> {{ countdown }}s </el-tag>
                            </span>
                            <span v-else class="mask-text">•••••••••••••••••••••••••••••••</span>
                        </div>
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('sysWallet.privateKey')">
                        <div class="sensitive-field" v-longpress="() => handleLongPress('privateKey')">
                            <span v-if="showPrivateKey" class="private-key-text"> {{ walletData.pri_key }} </span>
                            <span v-else class="mask-text">•••••••••••••••••••••••••••••••</span>
                        </div>
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('sysWallet.address')">{{ walletData.ac }}</el-descriptions-item>
                    <el-descriptions-item :label="t('sysWallet.path')">{{ walletData.path }}</el-descriptions-item>
                    <el-descriptions-item :label="t('sysWallet.currentIndex')">{{ walletData.current_index
                        }}</el-descriptions-item>
                    <el-descriptions-item :label="t('sysWallet.createdAt')"> {{
                        dayjs(walletData.created_at).format('YYYY-MM-DD HH:mm:ss') }} </el-descriptions-item>
                    <!-- <el-descriptions-item label="备注"> -->
                    <!-- <el-input v-model="remarkEdit" size="small" style="width: 200px" /> -->
                    <!-- <el-button type="primary" size="small" @click="updateRemark" class="ml-2">更新备注</el-button> -->
                    <!-- </el-descriptions-item> -->
                </el-descriptions>
            </div>
        </el-card>
        <!-- 创建确认对话框 -->
        <el-dialog :title="t('sysWallet.createTitle')" v-model="showCreateDialog" width="500px">
            <div class="create-warning">
                <el-alert type="error" :closable="false">
                    <template #title>
                        <div class="warning-content">
                            <el-icon><warning-filled /></el-icon>
                            <div>
                                <p>{{ t('sysWallet.createWarning') }}</p>
                                <ul>
                                    <li>• {{ t('sysWallet.createWarning1') }}</li>
                                    <li>• {{ t('sysWallet.createWarning2') }}</li>
                                    <li>• {{ t('sysWallet.createWarning3') }}</li>
                                    <li>• {{ t('sysWallet.createWarning4') }}</li>
                                </ul>
                            </div>
                        </div>
                    </template>
                </el-alert>
            </div>
            <template #footer>
                <el-button @click="showCreateDialog = false">{{ t('common.cancel') }}</el-button>
                <el-button type="danger" @click="confirmCreate">{{ t('common.confirm') }}</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup>
import { ref, onMounted, reactive } from 'vue'
import api from '../../../api'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { WarningFilled } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 在setup作用域顶部添加处理函数
const handleLongPress = (type) => {
    if (type === 'mnemonic') {
        showMnemonic.value = true
    } else {
        showPrivateKey.value = true
    }
    startCountdown()
}

// 修改长按指令逻辑
const vLongpress = {
    mounted(el, binding) {
        let pressTimer = null
        let autoHideTimer = null

        const start = () => {
            pressTimer = setTimeout(() => {
                binding.value() // 触发显示
                // 15秒后自动隐藏
                autoHideTimer = setTimeout(() => {
                    showMnemonic.value = false
                    showPrivateKey.value = false
                    resetCountdown()
                }, 15000)
            }, 1000)
        }

        const cancel = () => {
            // 仅取消长按触发，不取消自动隐藏
            clearTimeout(pressTimer)
        }

        // 事件监听（移除mouseup/touchend的隐藏事件）
        el.addEventListener('mousedown', start)
        el.addEventListener('touchstart', start)
        el.addEventListener('mouseup', cancel)
        el.addEventListener('touchend', cancel)
        el.addEventListener('mouseleave', cancel)
    }
}

const walletData = ref(null)
const showCreateDialog = ref(false)
const remarkEdit = ref('')
const sensitiveVisible = reactive({
    mnemonic: false,
    privateKey: false
})
const showMnemonic = ref(false)
const showPrivateKey = ref(false)
const countdown = ref(0)
let countdownTimer = null

// 加载系统钱包数据
const loadData = async () => {
    try {
        const { data } = await api.get('/sys-wallet')
        walletData.value = data
        remarkEdit.value = data.remark || ''
    } catch (error) {
        if (error.response?.status !== 404) {
            ElMessage.error('暂无系统钱包')
        }
    }
}

// 创建系统钱包
const confirmCreate = async () => {
    try {
        await api.post('/sys-wallet/create')
        ElMessage.success('系统钱包创建成功')
        showCreateDialog.value = false
        await loadData()
    } catch (error) {
        ElMessage.error('创建失败: ' + error.message)
    }
}

// 更新备注
const updateRemark = async () => {
    try {
        await api.put('/sys-wallet/update', {
            remark: remarkEdit.value
        })
        ElMessage.success('备注更新成功')
    } catch (error) {
        ElMessage.error('更新失败: ' + error.message)
    }
}

// 处理创建
const handleCreate = () => {
    showCreateDialog.value = true
}

// 倒计时逻辑
const startCountdown = () => {
    countdown.value = 15
    countdownTimer = setInterval(() => {
        if (countdown.value > 0) {
            countdown.value--
        } else {
            resetCountdown()
        }
    }, 1000)
}

const resetCountdown = () => {
    clearInterval(countdownTimer)
    countdown.value = 0
    showMnemonic.value = false
    showPrivateKey.value = false
}

onMounted(() => {
    loadData()
})
</script>
<style scoped lang="scss">
.security-alert {
    margin-bottom: 20px;

    .security-title {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #f56c6c;

        .tip-text {
            font-size: 12px;
            color: #e6a23c;
            margin-top: 4px;
        }
    }
}

.sensitive-field {
    padding: 8px;
    background: #f5f7fa;
    border-radius: 4px;
    cursor: pointer;
    transition: background 0.3s;
    position: relative;

    &:hover {
        background: #ebeef5;
    }

    .mnemonic-text {
        word-break: break-all;
        display: inline-block;
        line-height: 1.8;
        color: #f56c6c;
    }

    .private-key-text {
        word-break: break-all;
        color: #e6a23c;
    }

    .mask-text {
        letter-spacing: 2px;
        color: #909399;
    }
}

.warning-content {
    display: flex;
    gap: 12px;
    align-items: flex-start;

    ul {
        margin: 8px 0 0 20px;

        li {
            line-height: 1.8;
            color: #f56c6c;
        }
    }
}

.ml-2 {
    margin-left: 8px;
}

.countdown-tag {
    position: absolute;
    right: 8px;
    top: 8px;
    font-size: 12px;
    padding: 0 5px;
}
</style>