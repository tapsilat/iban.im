<template>
        <div class="bg-white rounded-lg p-4 shadow">
                <div class="flex flex-wrap gap-2 justify-center mb-4">
                        <router-link
                                v-for="tab in tabs"
                                :key="tab.id"
                                :to="tab.route"
                                exact
                                class="px-3 py-2 rounded hover:bg-gray-100"
                                active-class="bg-teal-600 text-white"
                        >
                                {{ tab.name }}
                        </router-link>
                </div>
                <div>
                        <router-view />
                </div>
        </div>
</template>

<script>
    export default {
        name: "Dashboard",
                                data() {
                                        return {
                                                tabs: [
                                                        { id: 1, name: "Profile", route: `/dashboard` },
                                                        { id: 2, name: "Security", route: `/dashboard/security` },
                                                        { id: 3, name: "IBANs", route: `/dashboard/ibans` },
                                                        { id: 4, name: "Logout", route: `/dashboard/logout` },
                                                ]
                                        }
                                }
    }
</script>