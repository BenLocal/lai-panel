<script setup lang="ts">
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import { setToastService } from '@/lib/toast'

const toast = useToast()
setToastService(toast)
</script>

<template>
    <router-view />
    <Toast />
</template>

<style scoped></style>
