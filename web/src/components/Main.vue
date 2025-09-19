<template>
    <div class="box">
         Shorten, create and share memorable links for IBANS.

        <strong>IBAN.im/nickname/bank</strong>
        <br><br>
        <router-link class="btn"  style="--color: var(--primary-2)" to="/login">
            Login
        </router-link>
        or
        <router-link class="btn" to="/register" >
            Register
        </router-link>

        Build By <a href="https://tapsilat.com">tapsilat</a>
        
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