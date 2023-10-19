<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const name = ref('')

const handleSubmit = () => {
  router.push({ name: 'room', query: { name: name.value } })
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <article>
      <header>Join a room</header>
      <input v-model="name" type="text" name="name" placeholder="Your Name" />
      <footer>
        <input type="submit" value="Join" />
      </footer>
    </article>
  </form>
</template>
