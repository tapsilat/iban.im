<template>
        <div v-if="profile" class="">
        <div class="single-iban__button mb-4">
            <router-link class="btn" :to="{ name: 'home.single', params: { username: $route.params.username } }">Profile</router-link>
        </div>
            <ul v-if="current && !current.isPrivate" class="list-none flex flex-col items-center">
                <li class="w-full max-w-xl flex border border-gray-200 border-b-0"> <span class="w-40 font-medium p-2">IBAN name</span><span class="p-2">{{ name }}</span></li>
                <li class="w-full max-w-xl flex border border-gray-200 border-b-0"> <span class="w-40 font-medium p-2">Handle</span><span class="p-2">{{ current.handle }}</span></li>
                <li class="w-full max-w-xl flex border border-gray-200 border-b-0">
                    <span class="w-40 font-medium p-2">IBAN</span>
                    <span class="p-2 flex-1">{{ current.text }}</span>
                    <span class="p-2"><button class="px-3 py-1 rounded bg-gray-100" @click="copy(current.text)">Copy</button></span>
                </li>
                <li class="w-full max-w-xl flex border border-gray-200">
                    <span class="p-2">{{ current.description }}</span>
                </li>
            </ul>
        <form v-else-if="current && current.isPrivate" class="show-info" @submit.prevent="showInfo">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                      <input v-model="formData.password" type="password" placeholder="Password" class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" required />
                </div>
                <div class="show-submit">
                    <button type="submit" class="btn mono-bg text-white px-4 py-2 rounded">Show</button>
                </div>
            </div>
        </form>
        <div v-else>
            <b>An account named <i>{{ $route.params.alias }}</i> was not found</b>
        </div>
    </div>
</template>

<script>
    import {mapActions, mapState} from 'vuex';
    export default {
        name: "SingleIban",
        data: () => ({
            canShow: false,
            formData: {
                password: '',
            },
            
        }),
        computed: {
            ...mapState(['ibans','profile']),
            current() {
                return this.ibans.filter( iban => iban.handle.toLowerCase() === this.$route.params.alias.toLowerCase())[0]
            },
            name() {
                return `${this.profile.firstName} ${this.profile.lastName}`
            }
        },
        created() {
            this.fetchSingleProfile({
                username : this.$route.params.username
            });
        },
        methods: {
            ...mapActions({
                fetchSingleProfile: 'fetchSingleProfile',
                checkShowPassword: 'checkShowPassword',
            }),
            showInfo() {
                this.checkShowPassword({
                    id       : this.current.id,
                    password : this.formData.password
                });
            },
            async copy(text) {
                try {
                    await navigator.clipboard.writeText(text);
                    alert('Iban was copied to clipboard!');
                } catch (e) {
                    alert('Something went wrong!');
                }
            }
        },
        watch : {
            '$store.state.canShow'(ok) {
                if(this.current.isPrivate && ok) {
                    this.current.isPrivate = false;
                }
            },
        }
    }
</script>

<style>

</style>