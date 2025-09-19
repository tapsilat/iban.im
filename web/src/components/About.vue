<template>
    <div class="box">
        <h1>About IBAN.im</h1>
        <p>I feel that I could never memorize numbers.<br />

        If you need to share it, it becomes a huge pain. You need to login your bank and copy if you can. Then share it via sms or WhatsApp. Then wait for next pain. Bank applications are hard when you have to copy/paste something.</p>
        <p>This service is build by <a href="https://tapsilat.com">tapsilat</a>.</p>
    </div>
</template>

<script>

    export default {
        name: "About",
        data: () => ({

        }),
        created() {
            this.$store.state.isLoaded = true
        }
    }
</script>