<template>
    <div class="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow">
        <h3 class="text-center text-lg font-semibold mb-4">Update Password</h3>
        <form class="grid grid-cols-1 md:grid-cols-2 gap-4" @submit.prevent="submit">
            <div>
                <label class="block text-sm font-medium">Password</label>
                <input
                    :type="showPassword ? 'text' : 'password'"
                    v-model="password"
                    class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                    minlength="7"
                />
                <button type="button" class="text-xs text-blue-600 mt-1" @click="showPassword = !showPassword">
                    {{ showPassword ? 'Hide' : 'Show' }}
                </button>
            </div>
            <div>
                <label class="block text-sm font-medium">Password Again</label>
                <input
                    :type="showPassword ? 'text' : 'password'"
                    v-model="passwordRepeat"
                    class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                    minlength="7"
                />
            </div>
            <div class="md:col-span-2 flex justify-end">
                <button
                    type="submit"
                    :disabled="!isValid"
                    class="px-4 py-2 rounded bg-blue-600 text-white disabled:bg-gray-300"
                >
                    Save
                </button>
            </div>
        </form>
    </div>
</template>

<script>
    export default {
        name: "Security",
        data: () => ({
            showPassword: false,
            passwordRepeat: null,
            password: null,
        }),
        computed: {
            isValid() {
                return (
                  this.password &&
                  this.passwordRepeat &&
                  this.password.length >= 7 &&
                  this.password === this.passwordRepeat
                );
            },
        },
        methods: {
            submit() {
                console.log('submitted');
                this.$store.dispatch('changePassword', {
                    password: this.password
                })
            }
        }
    }
</script>
