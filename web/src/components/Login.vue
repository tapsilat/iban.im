<template>
    <div class="box max-w-md mx-auto bg-white p-6 rounded-lg shadow">
        <h2 class="text-center text-xl font-semibold mb-4">Login</h2>
        <form class="space-y-4" @submit.prevent="submit">
            <div>
                <label class="block text-sm font-medium">Email</label>
            <input v-model="formData.email" type="email" class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" required />
            </div>
            <div>
                <label class="block text-sm font-medium">Password</label>
            <input v-model="formData.password" type="password" class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" required />
            </div>
            <div v-if="error" class="error">{{ error }}</div>
            <div class="flex justify-between">
                <router-link class="btn" to="/register">Register</router-link>
                <button type="submit" class="btn mono-bg text-white px-4 py-2 rounded">Login</button>
            </div>
        </form>
    </div>
</template>

<script>

    import { mapState,mapActions } from 'vuex';

    export default {
        name: "Login",
        data: () => ({
            formData: {
                email: '',
                password: '',
            },
        }),
        computed: {
            ...mapState(['error'])
        },
        created() {
            this.$store.dispatch('resetError');
            this.setLoaded(true);
        },
        methods: {
            ...mapActions({
                setLoaded: 'setLoaded',
            }),
            submit() {
                this.$store.dispatch('login', this.formData);
            }
        }
    }
</script>