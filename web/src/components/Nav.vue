<template>
    <div id="nav" class="container mx-auto p-4 flex justify-center items-center">
        <router-link class="home-logo" :to="this.$store.state.logged ? '/dashboard' : '/'" >
            IBAN.im
        </router-link>
    </div>
</template>

<script>
    export default {
        name: "Nav"
    }
</script>

<style scoped>

</style>