export default {
    common: {
        logout: 'Logout',
        confirm: 'Confirm',
        cancel: 'Cancel',
        edit: 'Edit',
        delete: 'Delete',
        create: 'Create',
        search: 'Search',
        details: 'Details',
        actions: 'Actions',
        id: 'ID',
        merchantName: 'Merchant Name',
        email: 'Email',
        phone: 'Phone',
        status: 'Status',
        createdAt: 'Created At',
        updatedAt: 'Updated At',
        enabled: 'Enabled',
        disabled: 'Disabled',
        enable: 'Enable',
        disable: 'Disable',
        confirmTitle: 'Confirmation',
        confirmDelete: 'Confirm Delete',
        operationFailed: 'Operation failed',
        update: 'Update',
        cancel: 'Cancel',
        remark: 'Remark',
        notSet: 'Not set',
        home: 'Home',
        profile: 'Profile',
        retry: 'Retry',
        confirmPayment: 'Confirm Payment',
        backHome: 'Back to Home',
        language: 'Language',
        zh: 'Chinese',
        en: 'English',
        notify_status: 'Notify Status'
    },
    login: {
        title: 'Welcome Back 👋',
        subtitle: 'Please sign in to your account',
        email: 'Email',
        password: 'Password',
        submit: 'Sign In'
    },
    dashboard: {
        dashboard: 'Dashboard',
        merchantTotal: 'Total Merchants',
        orderTotal: 'Total Orders',
        success7Days: '7-day Success',
        totalAmount: 'Total Amount',
        welcome: 'If you like cheems happy, please support us. pay.cheemshappy.com',
        platform: 'Payment Platform'
    },
    order: {
        merchant: 'Merchant',
        orders: 'Orders',
        searchOrder: 'Search Order No.',
        orderNo: 'Order No.',
        merchantName: 'Name',
        amount: 'Amount',
        createdAt: 'Created At',
        detailTitle: 'Order Details',
        merchantId: 'Merchant ID',
        systemOrderNo: 'System Order No.',
        merchantOrderNo: 'Merchant Order No.',
        currentStatus: 'Current Status',
        callbackUrl: 'Callback URL',
        updatedAt: 'Last Updated',
        deletedAt: 'Deleted At',
        paymentOrder: 'Payment Order',
        targetNetwork: 'Target Network',
        paymentSuccess: 'Payment Successful!',
        autoRedirect: 'seconds auto redirect to merchant page',
        jumpNow: 'Jump Now'
    },
    status: {
        success: 'success',
        pending: 'warning',
        failed: 'danger',
        canceled: 'info'
    },
    statusText: {
        completed: 'Completed',
        created: 'Created',
        failed: 'Failed',
        canceled: 'Canceled',
        listening: 'Listening',
        unknown: 'Unknown Status'
    },
    components: {
        pagination: {
            layout: 'total, sizes, prev, pager, next, jumper',
            total: 'Total {total} items'
        }
    },
    merchants: {
        merchants: 'Merchants',
        create: 'Create Merchant',
        editTitle: 'Edit Merchant',
        createTitle: 'Create Merchant',
        merchantName: 'Name',
        searchPlaceholder: 'Search merchant name',
        walletManagement: 'Wallet Management',
        rules: {
            nicknameRequired: 'Please enter merchant name',
            emailRequired: 'Please enter email address',
            emailFormat: 'Please enter valid email format',
            phoneFormat: 'Please enter valid phone number'
        },
        updateSuccess: 'Merchant updated successfully',
        createSuccess: 'Merchant created successfully',
        deleteSuccess: 'Merchant deleted successfully',
        deleteFailed: 'Delete failed'
    },
    wallet: {
        wallets: 'Wallets',
        updateTitle: 'Update Wallet',
        createTitle: 'Create Wallet',
        address: 'Address',
        privateKey: 'Private Key',
        mnemonic: 'Mnemonic',
        autoGenerate: 'Auto Generate',
        updateSuccess: 'Wallet updated successfully',
        createSuccess: 'Wallet created successfully',
        autoGenerateFailed: 'Auto generate failed',
        searchPlaceholder: 'Search wallet address',
        merchant: 'Merchant',
        balance: 'Balance',
        detailTitle: 'Wallet Details',
        autoGenerateTitle: 'Auto Generate Wallets',
        selectMerchant: 'Select Merchant',
        generateCount: 'Generate Count',
        startGenerate: 'Start Generate',
        generateSuccess: 'Successfully generated {count} wallets',
        generateFailed: 'Generation failed',
        statusUpdateSuccess: 'Status updated successfully',
        deleteConfirm: 'Confirm to delete this wallet? This action cannot be undone!',
        disableConfirm: 'Confirm to disable this wallet?',
        enableConfirm: 'Confirm to enable this wallet?',
        path: 'Path',
        create: 'Create Wallet',
        editWallet: 'Edit Wallet',
        deleteWallet: 'Delete Wallet',
        walletAddress: 'Address',
        privateKey: 'Private Key',
        mnemonic: 'Mnemonic',
        path: 'Path',
        connectMetaMask: 'Connect MetaMask',
        connected: 'Wallet Connected',
        ensureNetwork: 'Please ensure using a wallet that supports {network}'
    },
    settings: {
        systemWallet: 'System Wallet',
        feeConfig: 'Fee Configuration',
        notification: 'Notifications'
    },
    sysWallet: {
        mainAddress: 'Main Address',
        balance: 'Balance',
        lastSync: 'Last Sync',
        importWallet: 'Import Wallet',
        confirmImport: 'Confirm Import',
        privateKey: 'Private Key',
        rules: {
            privateKeyRequired: 'Please enter private key'
        },
        mnemonic: 'Mnemonic',
        currentIndex: 'Current Index',
        createdAt: 'Created At',
        address: 'Address',
        path: 'Path',
        updateSuccess: 'Wallet updated successfully',
        updateFailed: 'Update failed: ',
        importSuccess: 'Wallet imported successfully',
        importFailed: 'Import failed: ',
        securityWarning: 'Private key and mnemonic are important information for the system wallet, please keep them safe',
        longPressTip: 'Long press to view sensitive information, hidden after 15 seconds',
        noWallet: 'No wallet created',
        createWarning: 'About to create a new system wallet, please note:',
        createWarning1: 'This action will generate a new Ethereum wallet',
        createWarning2: 'Please immediately back up the mnemonic and private key',
        createWarning3: 'This operation cannot be modified after creation',
        createWarning4: 'Ensure the operation is performed in a secure environment',
        walletManagement: 'System Wallet Management'
    },
    otherSettings: {
        otherSettings: 'Other Settings',
        featureUnderDevelopment: 'Feature under development'
    },
    systemSettings: {
        systemSettings: 'System Settings',
        domain: 'Domain',
        domainPlaceholder: 'Please enter domain name',
    },
    login: {
        login: 'Login',
        welcome: 'Welcome Back 👋',
        pleaseLogin: 'Please sign in to your account'
    },
    merchantsApi: {
        merchantsApi: 'API Management',
        apikey: 'API Key',
        callbackUrl: 'Callback URL',
        secretKey: 'Secret Key',
        create: 'Create API',
        editTitle: 'Edit API Config',
        createTitle: 'Create API Config',
        searchPlaceholder: 'Search merchant/API key',
        deleteConfirm: 'Confirm to delete this API config? This action cannot be undone!',
        merchant: 'Merchant'
    },
    rules: {
        required: 'Required',
        invalidDomain: 'Invalid domain format'
    },
    error: {
        orderNotFound: 'Order Not Found',
        orderExpired: 'This order may have been paid or expired',
        connectFailed: 'Wallet Connection Failed',
        paymentFailed: 'Payment Failed'
    }
    // 其他翻译项...
} 