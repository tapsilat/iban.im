<template>
  <div>
    <Nav />
    <div class="container mx-auto p-4">
      <router-view />
    </div>
  </div>
</template>

<script>
import Nav from "./components/Nav.vue";

export default {
  name: "App",
  components: {
    Nav,
  },

  data: () => ({
    //
  }),
};
</script>
