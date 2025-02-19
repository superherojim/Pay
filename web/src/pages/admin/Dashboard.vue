<template>
    <div class="dashboard">
        <!-- 统计卡片 -->
        <el-skeleton :loading="loading" animated :count="4">
            <template #template>
                <el-skeleton-item variant="rect" style="width: 23%; height: 120px; margin: 0 1%; border-radius: 8px" />
            </template>
            <template #default>
                <div class="stats-grid">
                    <el-card v-for="stat in stats" :key="stat.title" class="stat-card" shadow="hover">
                        <div class="card-content">
                            <i :class="stat.icon" :style="{ color: stat.color }"></i>
                            <div class="text">
                                <div class="title">{{ stat.title }}</div>
                                <div class="value">{{ formatValue(stat) }}</div>
                            </div>
                        </div>
                    </el-card>
                </div>
            </template>
        </el-skeleton>
        <!-- 图表区域 -->
        <el-card class="chart-card" shadow="never">
            <template #header>
                <div class="chart-header">
                    <span style="font-size: 18px;font-weight: 500">{{ t('dashboard.welcome') }}</span>
                </div>
            </template>
            <div class="chart-container">
                <!-- 这里可以接入ECharts或其它图表库 -->
                <div class="chart-placeholder">
                    <i class="el-icon-data-analysis"></i>
                    <p>{{ t('dashboard.platform') }}</p>
                </div>
            </div>
        </el-card>
    </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { statsAPI } from '../../api'
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const stats = ref([
    {
        title: t('dashboard.merchantTotal'),
        value: 0,
        icon: 'el-icon-shopping-bag-1',
        color: '#409EFF',
        key: 'merchantTotal'
    },
    {
        title: t('dashboard.orderTotal'),
        value: 0,
        icon: 'el-icon-document',
        color: '#67C23A',
        key: 'orderTotal'
    },
    {
        title: t('dashboard.success7Days'),
        value: 0,
        icon: 'el-icon-success',
        color: '#E6A23C',
        key: 'success7Days'
    },
    {
        title: t('dashboard.totalAmount'),
        value: '¥ 0.00',
        icon: 'el-icon-money',
        color: '#F56C6C',
        key: 'orderTotal'
    }
])

const chartType = ref('week')
const loading = ref(true)

// 新增格式化方法
const formatValue = (stat) => {
    if (stat.key === 'totalAmount') return stat.value
    return stat.value.toLocaleString() // 数字添加千位分隔符
}

onMounted(async () => {
    try {
        const { data } = await statsAPI.getDashboardStats()

        stats.value = stats.value.map(item => ({
            ...item,
            value: item.title === '总交易额'
                ? `¥ ${(data.orderTotal * 100).toLocaleString()}`
                : data[item.key] || 0
        }))

    } catch (error) {
        console.error('获取统计失败:', error)
    } finally {
        loading.value = false
    }
})
</script>
<style scoped lang="scss">
.dashboard {
    max-width: 1200px;
    margin: 0 auto;
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 20px;
    margin-bottom: 24px;
}

.stat-card {
    :deep(.el-card__body) {
        padding: 20px;
    }

    .card-content {
        display: flex;
        align-items: center;

        i {
            font-size: 36px;
            margin-right: 16px;
        }

        .title {
            color: #909399;
            font-size: 14px;
            margin-bottom: 8px;
        }

        .value {
            font-size: 24px;
            font-weight: 600;
            color: #303133;
        }
    }
}

.chart-card {
    margin-top: 24px;

    .chart-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0 12px;
    }

    .chart-container {
        height: 400px;
    }
}

.chart-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #909399;

    i {
        font-size: 80px;
        margin-bottom: 20px;
        opacity: 0.6;
    }
}

@media (max-width: 768px) {
    .stats-grid {
        grid-template-columns: 1fr;
    }

    .chart-container {
        height: 300px;
    }
}

/* 新增加载样式 */
:deep(.el-skeleton) {
    width: 100%;
}
</style>
