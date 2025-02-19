<template>
    <div class="merchants-management">
        <el-card class="toolbar-card">
            <div class="toolbar">
                <el-button type="primary" @click="handleCreate">
                    <i class="el-icon-plus"></i>{{ t('merchants.create') }}</el-button>
                <el-input v-model="searchKey" :placeholder="t('merchants.searchPlaceholder')"
                    style="width: 300px; margin-left: 15px" clearable @clear="loadData">
                    <template #append>
                        <el-button icon="el-icon-search" @click="loadData" />
                    </template>
                </el-input>
            </div>
        </el-card>
        <el-card>
            <el-table :data="merchantsList" v-loading="loading" style="width: 100%" :max-height="500">
                <el-table-column prop="id" :label="t('common.id')" width="80" />
                <el-table-column prop="nickname" :label="t('merchants.merchantName')" width="400" />
                <el-table-column prop="email" :label="t('common.email')" />
                <el-table-column prop="phone" :label="t('common.phone')" width="150" />
                <el-table-column :label="t('common.actions')" width="300">
                    <template #default="scope">
                        <div class="action-buttons">
                            <el-button size="small" @click="handleEdit(scope.row)" type="primary"> {{ t('common.edit')
                            }} </el-button>
                            <el-button size="small" type="danger" @click="handleDelete(scope.row.id)"> {{
                                t('common.delete') }} </el-button>
                            <el-button size="small" @click="handleWallet(scope.row)" type="success"> {{
                                t('merchants.walletManagement') }} </el-button>
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
        <!-- 创建/编辑对话框 -->
        <el-dialog :title="isEdit ? t('merchants.editTitle') : t('merchants.createTitle')" v-model="showCreateDialog"
            width="500px" @close="handleDialogClose">
            <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
                <el-form-item :label="t('merchants.merchantName')" prop="nickname">
                    <el-input v-model="form.nickname" />
                </el-form-item>
                <el-form-item :label="t('common.email')" prop="email">
                    <el-input v-model="form.email" type="email" />
                </el-form-item>
                <el-form-item :label="t('common.phone')" prop="phone">
                    <el-input v-model="form.phone" />
                </el-form-item>
                <el-form-item label="密码" prop="password" v-if="!isEdit">
                    <el-input v-model="form.password" type="password" show-password />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="showCreateDialog = false">取消</el-button>
                <el-button type="primary" @click="submitForm">{{ t('common.confirm') }}</el-button>
            </template>
        </el-dialog>
        <!-- 钱包管理对话框 -->
        <el-dialog :title="isWalletExist ? t('wallet.updateTitle') : t('wallet.createTitle')" v-model="walletDialog"
            width="500px">
            <el-form :model="walletForm" :rules="walletRules" ref="walletFormRef" label-width="100px">
                <el-form-item :label="t('wallet.address')">
                    <el-input v-model="walletForm.ac" />
                </el-form-item>
                <el-form-item :label="t('wallet.privateKey')">
                    <el-input v-model="walletForm.pri_key" show-password />
                </el-form-item>
                <el-form-item :label="t('wallet.mnemonic')">
                    <el-input v-model="walletForm.particle_device" />
                </el-form-item>
                <el-form-item label="钱包路径" prop="path">
                    <el-input v-model="walletForm.path" />
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="walletForm.remark" type="textarea" />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="walletDialog = false">{{ t('common.cancel') }}</el-button>
                <el-button type="primary" @click="submitWalletForm"> {{ isWalletExist ? t('common.update') :
                    t('common.create') }} </el-button>
                <el-button type="warning" @click="generateWallet">{{ t('wallet.autoGenerate') }}</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const merchantsList = ref([])
const loading = ref(false)
const searchKey = ref('')
const showCreateDialog = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const walletFormRef = ref(null)
const walletDialog = ref(false)
const isWalletExist = ref(false)
const walletForm = reactive({
    nickname: '',
    mid: null,
    ac: '',
    pri_key: '',
    particle_device: '',
    path: '',
    remark: ''
})
const handleSizeChange = (size) => {
    pagination.size = size
    loadData()
}
// 分页参数
const pagination = reactive({
    index: 1,
    size: 10
})
const total = ref(0)

// 初始表单数据
const initialForm = {
    id: null,
    nickname: '',
    email: '',
    phone: '',
    password: ''
}

// 表单数据
const form = reactive({ ...initialForm })

// 表单验证规则
const rules = reactive({
    nickname: [
        { required: true, message: t('merchants.rules.nicknameRequired'), trigger: 'blur' }
    ],
    email: [
        {
            required: true,
            message: t('merchants.rules.emailRequired'),
            trigger: 'blur'
        },
        {
            type: 'email',
            message: t('merchants.rules.emailFormat'),
            trigger: ['blur', 'change']
        }
    ],
    phone: [
        {
            pattern: /^1[3-9]\d{9}$/,
            message: t('merchants.rules.phoneFormat'),
            trigger: 'blur'
        }
    ]
})

// 钱包验证规则
const walletRules = {
    ac: [{ required: true, message: '请输入钱包地址', trigger: 'blur' }],
    pri_key: [{ required: true, message: '请输入私钥', trigger: 'blur' }],
    particle_device: [{ required: true, message: '请输入助记词', trigger: 'blur' }],
    path: [{ required: true, message: '请输入钱包路径', trigger: 'blur' }]
}

// 加载数据
const loadData = async () => {
    try {
        loading.value = true
        const { data } = await api.post('/merchants/list', {
            ...pagination,
            search: searchKey.value
        })
        merchantsList.value = data.data
        total.value = data.total
    } finally {
        loading.value = false
    }
}

// 提交表单
const submitForm = async () => {
    await formRef.value.validate()

    try {
        const request = isEdit.value
            ? api.post('/merchants/update', form)
            : api.post('/merchants/create', form)

        await request
        ElMessage.success(
            isEdit.value ? t('merchants.updateSuccess') : t('merchants.createSuccess')
        )
        showCreateDialog.value = false
        resetForm()
        loadData()
    } catch (error) {
        ElMessage.error(
            `${isEdit.value ? t('common.update') : t('common.create')}${t('common.fail')}: ${error.response?.data?.message || error.message}`
        )
    }
}

// 重置表单方法
const resetForm = () => {
    Object.assign(form, initialForm)
    isEdit.value = false
    if (formRef.value) {
        formRef.value.resetFields()
    }
}

// 对话框关闭时重置表单
const handleDialogClose = () => {
    resetForm()
}

// 编辑商户
const handleEdit = (row) => {
    isEdit.value = true
    Object.assign(form, row)
    showCreateDialog.value = true
}

// 删除商户
const handleDelete = async (id) => {
    try {
        await ElMessageBox.confirm('确认删除该商户？此操作不可恢复！', '警告', {
            type: 'warning',
            confirmButtonText: '确认删除',
            cancelButtonText: '取消',
            confirmButtonClass: 'el-button--danger'
        })

        await api.delete(`/merchants/delete/${id}`)
        ElMessage.success(t('merchants.deleteSuccess'))
        loadData()
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error(t('merchants.deleteFailed') + error.message)
        }
    }
}

// 创建商户
const handleCreate = () => {
    resetForm()
    showCreateDialog.value = true
}

// 钱包管理
const handleWallet = async (merchant) => {
    walletDialog.value = true
    isWalletExist.value = false
    Object.assign(walletForm, {
        mid: merchant.id,
        nickname: merchant.nickname,
        ac: '',
        pri_key: '',
        particle_device: '',
        path: '',
        remark: ''
    })

    try {
        const { data } = await api.get(`/wallet/${merchant.id}`)
        if (data) {
            isWalletExist.value = true
            Object.assign(walletForm, data)
        }
    } catch (error) {
        console.error(error)
        ElMessage.error('获取钱包信息失败')
    }
}

const submitWalletForm = async () => {
    try {
        await walletFormRef.value.validate()
        const request = isWalletExist.value
            ? api.post('/wallet/update', {
                id: walletForm.id,
                m_id: walletForm.mid,
                ac: walletForm.ac,
                pri_key: walletForm.pri_key,
                particle_device: walletForm.particle_device,
                path: walletForm.path,
                remark: walletForm.remark
            })
            : api.post('/wallet/add', {
                m_id: walletForm.mid,
                ac: walletForm.ac,
                pri_key: walletForm.pri_key,
                particle_device: walletForm.particle_device,
                path: walletForm.path,
                remark: walletForm.remark
            })

        await request
        ElMessage.success(
            isWalletExist.value ? t('wallet.updateSuccess') : t('wallet.createSuccess')
        )
        walletDialog.value = false
        loadData()
    } catch (error) {
        ElMessage.error(`${isWalletExist.value ? t('common.update') : t('common.create')}${t('common.fail')}: ${error.response?.data?.message || error.message}`)
    }
}

const generateWallet = async () => {
    try {
        const endpoint = isWalletExist.value ? '/wallet/regenerate' : '/wallet/create'
        const { data } = await api.post(endpoint, {
            m_id: walletForm.mid
        })
        Object.assign(walletForm, data)
        const { data: newData } = await api.get(`/wallet/${walletForm.mid}`)
        Object.assign(walletForm, newData)
        isWalletExist.value = true
        ElMessage.success(
            isWalletExist.value ? t('wallet.updateSuccess') : t('wallet.createSuccess')
        )
    } catch (error) {
        ElMessage.error(
            `${t('wallet.autoGenerateFailed')}: ${error.message}`
        )
    }
}

onMounted(loadData)
</script>
<style scoped lang="scss">
@import './MerchantsManagement.scss';

.merchants-management {
    padding: 24px;
    background: #f5f6fa;
    min-height: calc(100vh - 60px);
}

/* 其他独有样式... */
</style>