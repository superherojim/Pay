<template>
    <div class="other-settings">
        <el-card>
            <template #header>
                <div class="card-header">
                    <span>{{ t('systemSettings.systemSettings') }}</span>
                </div>
            </template>
            <div class="setting-content">
                <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
                    <el-form-item :label="t('systemSettings.domain')" prop="domain">
                        <el-input v-model="form.domain" :placeholder="t('systemSettings.domainPlaceholder')" clearable
                            style="width: 400px" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" :loading="loading" @click="submitForm"> {{ t('common.confirm') }}
                        </el-button>
                    </el-form-item>
                </el-form>
            </div>
        </el-card>
    </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../../../api'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

const form = reactive({
    domain: ''
})

const rules = reactive({
    domain: [
        { required: true, message: t('rules.required'), trigger: 'blur' },
        {
            pattern: /^(?:(?:https?|ftp):\/\/)?(?:www\.)?[a-z0-9-]+(?:\.[a-z]{2,}){1,2}(?:\/.*)?$/i,
            message: t('rules.invalidDomain'),
            trigger: 'blur'
        }
    ]
})

const loading = ref(false)
const formRef = ref(null)
const existingConfig = ref(null)

// 加载系统配置
const loadConfig = async () => {
    try {
        const { data } = await api.get('/sys-config')
        if (data) {
            existingConfig.value = data
            form.domain = data.domain || ''
        }
    } catch (error) {
        console.error('加载配置失败:', error)
    }
}

// 提交表单
const submitForm = async () => {
    try {
        await formRef.value.validate()
        loading.value = true

        if (existingConfig.value) {
            // 更新现有配置
            await api.post('/sys-config/update', {
                ...existingConfig.value,
                domain: form.domain
            })
        } else {
            // 创建新配置
            await api.post('/sys-config/create', {
                domain: form.domain
            })
        }

        ElMessage.success(t('common.operationSuccess'))
        await loadConfig() // 重新加载最新配置
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error(t('common.operationFailed'))
        }
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadConfig()
})
</script>
<style scoped>
.setting-content {
    padding: 20px;
    max-width: 800px;
}
</style>