<template>
    <div class="wallet-management">
        <el-card class="toolbar-card">
            <div class="toolbar">
                <div>
                    <el-button type="primary" @click="handleCreate">
                        <i class="el-icon-plus"></i>{{ t('wallet.create') }}</el-button>
                    <el-button type="success" @click="showAutoDialog = true">
                        <i class="el-icon-magic-stick"></i>{{ t('wallet.autoGenerate') }}</el-button>
                </div>
                <el-input v-model="searchKey" :placeholder="t('wallet.searchPlaceholder')" clearable @clear="loadData"
                    style="width: 300px">
                    <template #append>
                        <el-button icon="el-icon-search" @click="loadData" />
                    </template>
                </el-input>
            </div>
        </el-card>
        <el-card>
            <el-table :data="walletList" v-loading="loading" style="width: 100%" :max-height="500">
                <el-table-column prop="id" label="ID" width="100" />
                <el-table-column :label="t('wallet.merchant')" width="220">
                    <template #default="scope"> {{merchantOptions.find(m => m.value === scope.row.m_id)?.label ||
                        '未知商户'}} </template>
                </el-table-column>
                <el-table-column prop="ac" :label="t('wallet.address')" width="500" />
                <el-table-column prop="path" :label="t('wallet.path')" />
                <!-- <el-table-column prop="balance" :label="t('wallet.balance')" width="120" align="right" /> -->
                <!-- <el-table-column :label="t('common.status')" width="100">
                    <template #default="scope">
                        <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" effect="dark"> {{ scope.row.status
                            === 1 ? t('common.enabled') : t('common.disabled') }} </el-tag>
                    </template>
                </el-table-column> -->
                <el-table-column :label="t('common.actions')" width="220" align="center">
                    <template #default="scope">
                        <div>
                            <el-button type="primary" link @click="handleEdit(scope.row)">{{ t('common.edit')
                            }}</el-button>
                            <el-button type="danger" link @click="handleDelete(scope.row.id)">{{ t('common.delete')
                            }}</el-button>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination">
                <el-pagination v-model:current-page="pagination.index" v-model:page-size="pagination.size"
                    :total="total" :layout="t('components.pagination.layout') || 'total, prev, pager, next'"
                    @current-change="loadData" @size-change="handleSizeChange" />
            </div>
        </el-card>
        <!-- 编辑对话框 -->
        <el-dialog :title="isEdit ? t('wallet.editTitle') : t('wallet.createTitle')" v-model="showDialog" width="500px">
            <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
                <el-form-item :label="t('wallet.merchant')" prop="m_id">
                    <el-select v-model="form.m_id" filterable :disabled="isEdit" :loading="loadingMerchants">
                        <el-option v-for="item in merchantOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('wallet.address')" prop="ac">
                    <el-input v-model="form.ac" />
                </el-form-item>
                <el-form-item :label="t('wallet.privateKey')" prop="pri_key">
                    <el-input v-model="form.pri_key" show-password />
                </el-form-item>
                <el-form-item :label="t('wallet.mnemonic')" prop="particle_device">
                    <el-input v-model="form.particle_device" />
                </el-form-item>
                <el-form-item :label="t('wallet.path')" prop="path">
                    <el-input v-model="form.path" />
                </el-form-item>
                <el-form-item :label="t('common.remark')" prop="remark">
                    <el-input v-model="form.remark" type="textarea" />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button type="primary" @click="submitForm"> {{ isEdit ? t('common.update') : t('common.create') }}
                </el-button>
                <el-button @click="showDialog = false">{{ t('common.cancel') }}</el-button>
            </template>
        </el-dialog>
        <!-- 自动生成对话框 -->
        <el-dialog :title="t('wallet.autoGenerateTitle')" v-model="showAutoDialog" width="500px">
            <el-form :model="autoForm" label-width="80px">
                <el-form-item :label="t('wallet.selectMerchant')" required>
                    <el-select v-model="autoForm.m_id" filterable :loading="loadingMerchants">
                        <el-option v-for="item in merchantOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="showAutoDialog = false">{{ t('common.cancel') }}</el-button>
                <el-button type="primary" @click="handleAutoCreate">{{ t('wallet.startGenerate') }}</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'

const { t } = useI18n()

const walletList = ref([])
const loading = ref(false)
const searchKey = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const pagination = reactive({
    index: 1,
    size: 10
})
const total = ref(0)

const form = reactive({
    id: null,
    m_id: null,
    ac: '',
    pri_key: '',
    particle_device: '',
    path: '',
    remark: ''
})

const rules = {
    m_id: [{
        required: true,
        message: '请选择商户',
        trigger: 'change'
    }],
    ac: [{ required: true, message: '请输入钱包地址', trigger: 'blur' }],
    pri_key: [{ required: true, message: '请输入私钥', trigger: 'blur' }],
    path: [{ required: true, message: '请输入路径', trigger: 'blur' }]
}

const merchantOptions = ref([])
const loadingMerchants = ref(false)

const loadMerchants = async () => {
    try {
        loadingMerchants.value = true
        const { data } = await api.get('/merchants/in')
        merchantOptions.value = data.map(m => ({
            value: m.id,
            label: m.nickname
        }))
    } finally {
        loadingMerchants.value = false
    }
}
const handleSizeChange = (size) => {
    pagination.size = size
    loadData()
}
const loadData = async () => {
    try {
        loading.value = true
        const { data } = await api.post('/wallet/list', {
            ...pagination,
            search: searchKey.value
        })
        walletList.value = data.data || []
        total.value = data.total || 0
    } catch (error) {
        ElMessage.error('加载数据失败')
    } finally {
        loading.value = false
    }
}

const handleEdit = (row) => {
    isEdit.value = true
    Object.assign(form, row)
    showDialog.value = true
}

const handleDelete = async (id) => {
    try {
        await ElMessageBox.confirm(t('wallet.deleteConfirm'), t('common.confirmTitle'), {
            type: 'error',
            confirmButtonText: t('common.confirmDelete'),
            cancelButtonText: t('common.cancel'),
            confirmButtonClass: 'el-button--danger'
        })
        await api.delete(`/wallet/delete/${id}`)
        ElMessage.success(t('common.deleteSuccess'))
        loadData()
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error(t('common.deleteFailed'))
        }
    }
}

const handleRegenerate = async (row) => {
    try {
        const { data } = await api.post('/wallet/regenerate', { m_id: row.m_id })
        ElMessage.success('重新生成成功')
        loadData()
    } catch (error) {
        ElMessage.error('生成失败: ' + error.message)
    }
}

const submitForm = async () => {
    try {
        await formRef.value.validate()

        if (isEdit.value) {
            await api.post('/wallet/update', form)
            ElMessage.success('更新成功')
        } else {
            // 普通手动创建仍需要验证
            await api.post('/wallet/add', form)
            ElMessage.success('创建成功')
        }

        showDialog.value = false
        loadData()
    } catch (error) {
        ElMessage.error(`提交失败: ${error.message}`)
    }
}

const generateWallet = async () => {
    try {
        const { data } = await api.post('/wallet/create', {
            m_id: form.m_id
        })

        // 合并生成的数据到表单
        Object.assign(form, {
            ...data,
            ac: data.ac || '',
            pri_key: data.pri_key || '',
            path: data.path || '/'
        })

        // 直接调用添加接口
        await api.post('/wallet/add', form)

        ElMessage.success('钱包生成并添加成功')
        showDialog.value = false
        loadData()  // 强制刷新列表

    } catch (error) {
        ElMessage.error(`生成失败: ${error.message}`)
    }
}

const handleCreate = () => {
    isEdit.value = false
    resetForm()
    showDialog.value = true
}

const resetForm = () => {
    Object.assign(form, {
        id: null,
        m_id: null,
        ac: '',
        pri_key: '',
        particle_device: '',
        path: '',
        remark: ''
    })
}

const showAutoDialog = ref(false)

const autoForm = reactive({
    m_id: null
})

const handleAutoCreate = async () => {
    try {
        const { data } = await api.post('/wallet/create', {
            m_id: autoForm.m_id
        })

        ElMessage.success('钱包生成成功')
        showAutoDialog.value = false
        loadData()
    } catch (error) {
        ElMessage.error(`生成失败: ${error.message}`)
    }
}

const handleAutoGenerate = async () => {
    try {
        await api.post('/wallet/batch-create', autoForm)
        ElMessage.success(t('wallet.generateSuccess', { count: autoForm.count }))
        showAutoDialog.value = false
        loadData()
    } catch (error) {
        ElMessage.error(`${t('wallet.generateFailed')}: ${error.message}`)
    }
}

const toggleStatus = async (wallet) => {
    try {
        await ElMessageBox.confirm(
            wallet.status === 1 ? t('wallet.disableConfirm') : t('wallet.enableConfirm'),
            t('common.confirmTitle'),
            {
                type: 'warning',
                confirmButtonText: t('common.confirm'),
                cancelButtonText: t('common.cancel')
            }
        )

        await api.post(`/wallet/status/${wallet.id}`, {
            status: wallet.status === 1 ? 0 : 1
        })
        ElMessage.success(t('wallet.statusUpdateSuccess'))
        loadData()
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error(t('common.operationFailed'))
        }
    }
}

onMounted(() => {
    loadData()
    loadMerchants()
})
</script>
<style scoped lang="scss">
/* 复用商户管理的样式 */
@import './MerchantsManagement.scss';

.toolbar {
    display: flex;
    align-items: center;
    gap: 12px;

    .el-button:nth-child(1) {
        margin-right: 12px;
    }

    .search-input {
        margin-left: auto;
    }
}
</style>