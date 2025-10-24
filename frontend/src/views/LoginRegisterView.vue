<template>
  <div class="login-register-view">
    <NavBar />
    <div class="login-register-container">
      <LoginRegister />
    </div>
  </div>
</template>

<script setup>
import { defineAsyncComponent } from 'vue'

const NavBar = defineAsyncComponent(() => import('@/components/NavBar.vue'))
const LoginRegister = defineAsyncComponent(() => import('@/components/LoginRegister.vue'))
</script>

<style scoped>
.login-register-view {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-width: 1200px;
  margin: auto;
}
</style>
