<template>
    <div v-if="profile">
        <div class="text-center mb-4">
            <h1 class="text-xl font-semibold">{{ name }}</h1>
        </div>
        <div v-if="ibans && ibans.length > 0">
            <ul class="space-y-2">
                <li v-for="(item,i) in ibans" :key="i">
                    <router-link class="inline-block px-3 py-2 rounded bg-gray-100 hover:bg-gray-200" :to="`/${profile.handle}/${item.handle}`">
                        {{ item.handle }}
                    </router-link>
                </li>
            </ul>
        </div>
        <div v-else class="text-center text-gray-500">
            No public IBANs available.
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex';
    export default {
        name: "Single",
        data: () => ({
        }),
        computed : {
            profile() {
                return this.$store.state.profile;
            },
            ibans() {
                return this.$store.state.ibans.filter(iban => !iban.isPrivate)
            },
            name() {
                return this.profile.firstName !== "" ? this.profile.firstName + " " + this.profile.lastName : this.profile.handle;
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
            }),
        }
    }
</script>

<style>

</style>