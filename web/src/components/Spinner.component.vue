<template>
  <span class="spinner material-symbols-outlined">progress_activity</span>
</template>

<style scoped lang="scss">
.spinner {
  font-size: 2rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
