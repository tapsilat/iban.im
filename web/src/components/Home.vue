<template>
    <div class="home-content bg-white rounded-lg p-4">
        <div class="flex flex-wrap gap-2 justify-center mb-4">
            <router-link
                v-for="tab in tabs"
                :key="tab.id"
                :to="tab.route"
                :class="['px-3 py-2 rounded hover:bg-gray-100', $route.path === tab.route ? 'bg-teal-600 text-white' : '']"
            >
                {{ tab.name }}
            </router-link>
        </div>
        <router-view class="i-tab" />
    </div>
  
</template>

<script>

    const publicRoutes = [
        { id: 1, name: "Home", route: `/`},
        { id: 2, name: "Login", route: `/login` },
        { id: 3, name: "Register", route: `/register` },
       { id: 4, name: "About", route: `/about`},
    ];

    const privateRoutes = [
        { id: 1, name: "Profile", route: `/dashboard`},
        { id: 2, name: "Security", route: `/dashboard/security` },
        { id: 3, name: "IBANs", route: `/dashboard/ibans` },
        { id: 4, name: "Scan", route: `/dashboard/scan` },
        { id: 5, name: "Logout", route: `/dashboard/logout` },
        { id: 6, name: "About", route: `/about` },
    ];

    export default {
        name: "Home",
        data() {
            return {
                activeTab: null,
                tabs: publicRoutes
            }
        },
        created() {
            if("username" in this.$route.params){
            }else{
                this.$store.state.isLoaded = true
            }
            if(this.$store.state.logged) {
                this.tabs = privateRoutes;
            }
        },
        watch : {
            '$store.state.logged'(value) {
                if(value) {
                    this.tabs = privateRoutes;
                }else{
                    this.tabs = publicRoutes;
                }
            },
        }
    }
</script>