<template>
  <div class="merchants-api-management">
    <el-card class="toolbar-card">
      <div class="toolbar">
        <div>
          <el-button type="primary" @click="handleCreate">
            <i class="el-icon-plus"></i>{{ t('merchantsApi.create') }} </el-button>
        </div>
        <el-input v-model="searchKey" :placeholder="t('merchantsApi.searchPlaceholder')" clearable @clear="loadData"
          style="width: 300px">
          <template #append>
            <el-button icon="el-icon-search" @click="loadData" />
          </template>
        </el-input>
      </div>
    </el-card>
    <el-card>
      <el-table :data="apiList" v-loading="loading" style="width: 100%" :max-height="500">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="m_id" label="MID" />
        <el-table-column :label="t('merchantsApi.merchant')" width="220">
          <template #default="scope"> {{merchantOptions.find(m => m.value === scope.row.m_id)?.label || '未知商户'}}
          </template>
        </el-table-column>
        <el-table-column prop="apikey" :label="t('merchantsApi.apikey')" />
        <el-table-column prop="callback_url" :label="t('merchantsApi.callbackUrl')" />
        <el-table-column :label="t('common.actions')" width="180" align="center">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">{{ t('common.edit') }}</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row.id)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination v-model:current-page="pagination.index" v-model:page-size="pagination.size" :total="total"
          :layout="t('components.pagination.layout') || 'total, prev, pager, next'" @current-change="loadData"
          @size-change="handleSizeChange" />
      </div>
    </el-card>
    <!-- 编辑对话框 -->
    <el-dialog :title="isEdit ? t('merchantsApi.editTitle') : t('merchantsApi.createTitle')" v-model="showDialog"
      width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item :label="t('merchantsApi.merchant')" prop="m_id">
          <el-select v-model="form.m_id" filterable :disabled="isEdit">
            <el-option v-for="item in merchantOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('merchantsApi.callbackUrl')" prop="callback_url">
          <el-input v-model="form.callback_url" />
        </el-form-item>
        <el-form-item :label="t('merchantsApi.secretKey')" prop="secret_key">
          <el-input v-model="form.secret_key" show-password />
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
  </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const { t } = useI18n()

// 数据列表
const apiList = ref([])
const loading = ref(false)
const searchKey = ref('')
const merchantOptions = ref([])

// 分页
const pagination = reactive({
  index: 1,
  size: 10
})
const total = ref(0)

// 表单相关
const showDialog = ref(false)
const isEdit = ref(false)
const form = reactive({
  m_id: null,
  callback_url: '',
  secret_key: '',
  remark: ''
})

const rules = {
  m_id: [{ required: true, message: t('rules.required'), trigger: 'blur' }],
  callback_url: [{ required: true, message: t('rules.required'), trigger: 'blur' }],
  secret_key: [{ required: true, message: t('rules.required'), trigger: 'blur' }]
}

// 加载数据
const loadData = async () => {
  try {
    loading.value = true
    const { data } = await api.post('/merchants/api/list', {
      page: pagination.index,
      size: pagination.size,
      keyword: searchKey.value
    })
    apiList.value = data.data
    total.value = data.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}
const handleSizeChange = (size) => {
  pagination.size = size
  loadData()
}
// 加载商户选项
const loadMerchants = async () => {
  try {
    const { data } = await api.get('/merchants/in')
    merchantOptions.value = data.map(m => ({
      value: m.id,
      label: m.nickname
    }))
  } catch (error) {
    console.error(error)
  }
}

// 表单操作
const handleCreate = () => {
  resetForm()
  isEdit.value = false
  showDialog.value = true
}

const handleEdit = (row) => {
  Object.assign(form, row)
  isEdit.value = true
  showDialog.value = true
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm(
      t('merchantsApi.deleteConfirm'),
      t('common.confirmTitle'),
      {
        type: 'warning',
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel')
      }
    )
    await api.delete(`/merchants/api/delete/${id}`)
    ElMessage.success(t('common.deleteSuccess'))
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('common.operationFailed'))
    }
  }
}

const submitForm = async () => {
  try {
    if (isEdit.value) {
      await api.post(`/merchants/api/update`, form)
    } else {
      await api.post('/merchants/api/create', form)
    }
    ElMessage.success(t('common.operationSuccess'))
    showDialog.value = false
    loadData()
  } catch (error) {
    console.error(error)
  }
}

const resetForm = () => {
  Object.assign(form, {
    m_id: null,
    callback_url: '',
    secret_key: '',
    remark: ''
  })
}

onMounted(() => {
  loadData()
  loadMerchants()
})
</script>
<style scoped lang="scss">
@import './MerchantsManagement.scss';
</style>