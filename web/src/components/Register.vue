<template>
    <div class="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow">
        <h2 class="text-center text-2xl font-semibold mb-6">Register</h2>

        <form class="grid grid-cols-1 md:grid-cols-2 gap-4" @submit.prevent="submit">
            <div>
                <label class="block text-sm font-medium mb-1">First Name</label>
                <input v-model="formData.firstName" class="w-full rounded border px-3 py-2" placeholder="First name" required />
            </div>
            <div>
                <label class="block text-sm font-medium mb-1">Last Name</label>
                <input v-model="formData.lastName" class="w-full rounded border px-3 py-2" placeholder="Last name" required />
            </div>
            <div>
                <label class="block text-sm font-medium mb-1">Username</label>
                <input v-model="formData.handle" class="w-full rounded border px-3 py-2 lowercase" placeholder="username" pattern="[A-Za-z0-9]*" required />
                <p class="text-xs text-gray-500 mt-1">Only letters and numbers.</p>
            </div>
            <div>
                <label class="block text-sm font-medium mb-1">Email</label>
                <input v-model="formData.email" type="email" class="w-full rounded border px-3 py-2" placeholder="you@example.com" required />
            </div>
            <div>
                <label class="block text-sm font-medium mb-1">Password</label>
                <div class="relative">
                    <input :type="showPassword ? 'text' : 'password'" v-model="formData.password" class="w-full rounded border px-3 py-2 pr-10" placeholder="Password" minlength="7" required />
                    <button type="button" class="absolute inset-y-0 right-0 px-3 text-sm text-gray-500" @click="showPassword = !showPassword">{{ showPassword ? 'Hide' : 'Show' }}</button>
                </div>
            </div>
            <div>
                <label class="block text-sm font-medium mb-1">Password Again</label>
                <input :type="showPassword ? 'text' : 'password'" v-model="passwordRepeat" class="w-full rounded border px-3 py-2" placeholder="Repeat password" required />
                <p v-if="passwordRepeat && !passwordsMatch" class="text-xs text-red-600 mt-1">Password does not match</p>
            </div>
            <div class="md:col-span-2 flex items-center gap-2">
                <input id="visible" type="checkbox" v-model="formData.visible" class="h-4 w-4" />
                <label for="visible" class="text-sm">Show my email address on my public profile.</label>
            </div>

            <div v-if="error" class="md:col-span-2 text-red-600">{{ error }}</div>

            <div class="md:col-span-2 flex justify-between pt-2">
                <router-link class="px-4 py-2 rounded border" to="/login">Login</router-link>
                <button :disabled="!formValid" class="px-4 py-2 rounded bg-blue-600 text-white disabled:opacity-50" type="submit">Register</button>
            </div>
        </form>
    </div>
  
</template>

<script>
    import { mapState } from 'vuex';

    export default {
        name: "Register",
        data: () => ({
            showPassword: false,
            formData: {
                firstName: '',
                lastName: '',
                email: '',
                handle: '',
                password: '',
                visible: false,
            },
            passwordRepeat: null,
        }),
        computed: {
            passwordsMatch() {
                return this.formData.password === this.passwordRepeat;
            },
            formValid() {
                return (
                  this.formData.firstName &&
                  this.formData.lastName &&
                  /.+@.+/.test(this.formData.email) &&
                  /^[A-Za-z0-9]*$/.test(this.formData.handle) &&
                  this.formData.password && this.formData.password.length >= 7 &&
                  this.passwordsMatch
                );
            },
            ...mapState(['error'])
        },
        methods: {
            submit() {
                if (!this.formValid) return;
                this.$store.dispatch('register', this.formData)
            }
        }
    }
</script>

<style scoped>
</style>