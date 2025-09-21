<template>
    <div class="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow text-center">
         <p>Shorten, create and share memorable links for IBANS.</p>

        <p><strong>IBAN.im/nickname/bank</strong></p>
        <p>
        <router-link class="inline-block px-4 py-2 rounded bg-blue-600 text-white" to="/login">Login</router-link>
        or
        <router-link class="inline-block px-4 py-2 rounded bg-gray-100" to="/register">Register</router-link>

        </p>
      
        <p>Build By <a href="https://tapsilat.com">tapsilat</a></p>
        
    </div>
</template>

<script>
    export default {
        name: "Main",
        data: () => ({

        }),
        created() {
            this.$store.state.isLoaded = true
        }
    }
</script>