<template>
    <div class="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow" v-if="model">
        <h3 class="text-center text-lg font-semibold mb-4">Profile</h3>
        <form class="grid grid-cols-1 md:grid-cols-2 gap-4" @submit.prevent="submit">
            <div class="md:col-span-2">
                <label class="block text-sm font-medium">Email</label>
                <input v-model="model.email" disabled class="w-full border border-gray-200 rounded px-3 py-2 bg-gray-50" />
            </div>
            <div>
                <label class="block text-sm font-medium">First Name</label>
                <input v-model="model.firstName" disabled class="w-full border border-gray-200 rounded px-3 py-2 bg-gray-50" />
            </div>
            <div>
                <label class="block text-sm font-medium">Last Name</label>
                <input v-model="model.lastName" disabled class="w-full border border-gray-200 rounded px-3 py-2 bg-gray-50" />
            </div>
            <div>
                <label class="block text-sm font-medium">Username</label>
                <input v-model="model.handle" class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div class="md:col-span-2">
                <label class="block text-sm font-medium">Bio</label>
                <textarea v-model="model.bio" rows="3" class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
            </div>
            <div v-if="error" class="md:col-span-2 text-red-600 text-sm">
                {{ error }}
            </div>
            <div class="md:col-span-2 flex justify-end">
                <button type="submit" class="px-4 py-2 rounded bg-blue-600 text-white">Save</button>
            </div>
        </form>
    </div>
</template>

<script>
    import { mapActions } from 'vuex';
    export default {
        name: "Profile",
        data: () => ({
            error: null,
            model: {
                email: '',
                firstName: '',
                lastName: '',
                handle: '',
                bio: '',
            },
            formRules: {
                handle: [
                    v => !!v || 'You need an username',
                    v => /^[A-Za-z0-9]*$/.test(v) || 'please only AZaz09'
                ],
            },
        }),
        created() {
            this.fetchProfile();
        },
        methods: {
            ...mapActions({
                fetchProfile: 'fetchProfile',
                changeProfile: 'changeProfile',
            }),
            submit() {
                console.log('submitted');
                this.error = null;
                this.changeProfile({
                    bio: this.model.bio,
                    handle: this.model.handle
                });
            }
        },
        watch: {
            '$store.state.profile'(value) {
                this.model = value;
            },
            '$store.state.error'(value) {
                this.error = value;
            }
        }
    }
</script>
