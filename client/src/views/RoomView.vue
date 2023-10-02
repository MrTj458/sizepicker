<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

let socket = null

onMounted(() => {
  socket = new WebSocket(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`)

  socket.addEventListener('message', (m) => {
    data.value = JSON.parse(m.data)
    console.log(data.value)
  })

  socket.addEventListener('error', () => {
    alert('Error creating websocket connection. Please refresh and try again')
  })

  socket.addEventListener('close', () => {
    alert('Websocket connection closed. Please refresh and try again')
  })

  socket.addEventListener('open', () => {
    socket.send(JSON.stringify({ cmd: 'name', name: route.query.name }))
  })
})

const vals = [1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144]
let data = ref({})

const choice = computed(() => {
  if (data.value.clients) {
    let c = data.value.clients?.find((c) => c.name === route.query.name)
    if (c) {
      return c.choice
    }
  }

  return 0
})

const pick = (num) => {
  console.log(num)
  socket.send(JSON.stringify({ cmd: 'pick', choice: num }))
}

const toggleShow = () => {
  if (data.value.show) {
    socket.send(JSON.stringify({ cmd: 'reset' }))
  } else {
    socket.send(JSON.stringify({ cmd: 'show' }))
  }
}
</script>

<template>
  <div class="main">
    <div class="cards mb">
      <article
        v-for="(client, idx) in data.clients"
        :key="idx"
        :class="{ selected: client.choice !== 0, 'selected-client': client.choice !== 0 }"
      >
        <header>
          {{ client.name }}
        </header>
        {{ data.show ? client.choice : '?' }}
      </article>
    </div>

    <div class="cards">
      <button @click="toggleShow" class="outline">{{ data.show ? 'Reset' : 'Show' }}</button>
    </div>

    <div class="cards pointer">
      <article
        v-for="val in vals"
        :key="val"
        @click="pick(val)"
        :class="{ selected: val === choice }"
      >
        {{ val }}
      </article>
    </div>
  </div>
</template>

<style scoped>
.main {
  margin-top: 3rem;
}

.cards {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.selected {
  color: var(--secondary-inverse);
  background-color: var(--secondary);
}

.selected-client header {
  background-color: var(--contrast);
  color: var(--contrast-inverse);
}

article {
  width: 100px;
  text-align: center;
  margin-bottom: 0;
  margin-top: 0;
}

.pointer {
  cursor: pointer;
}

.mb {
  margin-bottom: 1rem;
}

button {
  display: inline;
  width: auto;
}
</style>
