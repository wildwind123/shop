<template>
  <h1>You did it!</h1>
  <p>
    Visit <a href="https://vuejs.org/" target="_blank" rel="noopener">vuejs.org</a> to read the
    documentation
  </p>
  <div>
    <ButtonPrime :unstyled="false" :size="'small'">asdf</ButtonPrime>
  </div>
</template>
<script setup lang="ts">
import ButtonPrime from 'primevue/button'
</script>
