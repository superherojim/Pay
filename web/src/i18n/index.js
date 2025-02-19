import { createI18n } from 'vue-i18n'
import enLocale from './en'
import zhLocale from './zh'
import elementEnLocale from 'element-plus/es/locale/lang/en'
import elementZhLocale from 'element-plus/es/locale/lang/zh-cn'

const messages = {
    en: {
        ...enLocale,
        ...elementEnLocale
    },
    zh: {
        ...zhLocale,
        ...elementZhLocale
    }
}

const i18n = createI18n({
    legacy: false,
    locale: localStorage.getItem('lang') || 'zh',
    fallbackLocale: 'en',
    messages
})

export default i18n 