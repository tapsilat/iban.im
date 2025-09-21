<template>
    <div class="max-w-3xl mx-auto">
        <h3 class="text-center text-xl font-semibold mb-4">IBANs</h3>
        <div class="border rounded-md divide-y">
            <template v-if="ibans && ibans.length">
                <button
                    v-for="(item, i) in ibans"
                    :key="i"
                    class="w-full flex items-center justify-between px-4 py-3 hover:bg-gray-50 text-left lowercase"
                    @click="selectedIndex = i"
                >
                    <span class="font-medium">{{ item.handle }}</span>
                    <span class="text-gray-400">{{ selectedIndex === i ? '-' : '+' }}</span>
                </button>
            </template>
            <div v-else class="p-4 text-gray-500">No IBANs yet.</div>
        </div>

        <div class="mt-4">
            <button @click="show" class="inline-flex items-center gap-2 px-4 py-2 rounded-md bg-blue-600 text-white">
                <span class="text-lg">+</span>
                Add
            </button>
        </div>

        <!-- Modal -->
        <div v-if="dialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
            <div class="w-full max-w-xl bg-white rounded-lg shadow p-6 relative">
                <button @click="dialog = false" class="absolute right-3 top-3 text-gray-500 hover:text-gray-700" aria-label="Close">×</button>

                <h4 class="text-lg font-semibold mb-4">{{ current.id === '' ? 'Add IBAN' : 'Edit IBAN' }}</h4>

                <form @submit.prevent="save" class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium mb-1">Handle</label>
                        <input
                            v-model="current.handle"
                            class="w-full rounded border px-3 py-2 lowercase"
                            placeholder="handle"
                            required
                            pattern="[A-Za-z0-9]*"
                        />
                        <p class="text-xs text-gray-500 mt-1">Only letters and numbers.</p>
                    </div>

                    <div>
                        <label class="block text-sm font-medium mb-1">IBAN No</label>
                        <input v-model="current.text" class="w-full rounded border px-3 py-2" placeholder="TRXXXXXXXXXXXXXXXXXXXX" required />
                    </div>

                    <div>
                        <label class="block text-sm font-medium mb-1">IBAN Description</label>
                        <input v-model="current.description" class="w-full rounded border px-3 py-2" placeholder="Description" />
                    </div>

                    <div class="flex items-center gap-2">
                        <input id="isPrivate" type="checkbox" v-model="current.isPrivate" class="h-4 w-4" />
                        <label for="isPrivate" class="text-sm">Private</label>
                    </div>

                    <div v-if="current.isPrivate">
                        <label class="block text-sm font-medium mb-1">Password</label>
                        <input v-model="current.password" type="password" class="w-full rounded border px-3 py-2" placeholder="Password" required />
                    </div>

                    <div class="flex items-center justify-between pt-2">
                        <button v-if="current.id === ''" type="button" @click="cancel" class="px-4 py-2 rounded border">Cancel</button>
                        <button v-else type="button" @click="remove" class="px-4 py-2 rounded border border-red-300 text-red-700">Delete</button>
                        <div class="flex-1"></div>
                        <button type="submit" class="px-4 py-2 rounded bg-blue-600 text-white">Save</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapActions,mapState } from 'vuex';
    const cloneDeep = (obj) => JSON.parse(JSON.stringify(obj));

    function reset() {
        return {
            id: "",
            handle: '',
            text: '',
            description: '',
            isPrivate: false,
            password: '',
        }
    }

    export default {
        name: "Ibans",
        data: () => ({
            dialog: false,
            showForm: false,
            selectedIndex: undefined,
            current: reset(),
        }),
        computed: {
            ...mapState(['ibans']),
            passwordRule() {
                return () => (this.current.isPrivate && this.current.password !== '') || 'Please provide password'
            },
        },
        created() {
            this.fetchProfile();
            this.fetchIbans();
        },
        methods: {
            ...mapActions({
                fetchProfile: 'fetchProfile',
                fetchIbans: 'fetchIbans',
                ibanUpdate: 'ibanUpdate',
                ibanDelete: 'ibanDelete',
            }),
            remove() {
                console.log(this.selectedIndex);
                const r = confirm('Are you sure?');
                if (r !== true) {
                    return;
                }

                this.ibanDelete(this.ibans[this.selectedIndex].id).then((data) => {
                    if(data.errors){
                        alert(data.errors[0].message);
                        return;
                    }
                    if(data.data.ibanDelete.ok){
                        this.ibans.splice(this.selectedIndex, 1);
                        this.selectedIndex = undefined;
                        this.dialog = false;
                        this.current = reset();
                    }else{
                        alert(data.data.ibanDelete.msg);
                    }
                })
            },
            cancel() {
                this.dialog = false;
                this.selectedIndex = undefined;
                this.current = reset();
            },
            save() {
                // yoksa null hatası veriyor
                if(!this.current.isPrivate){
                    this.current.password = '';
                }
                const process = this.current.id === "" ? "ibanNew" : "ibanUpdate";
                this.ibanUpdate(this.current).then((data) => {
                    if(data.errors){
                        alert(data.errors[0].message);
                        return;
                    }
                    if(!data.data[process].ok){
                        alert(data.data[process].msg);
                    }else{
                        this.current.id = data.data[process].iban.id;
                        if(this.selectedIndex !== undefined) {
                            this.ibans[this.selectedIndex] = cloneDeep(this.current)
                        }else{
                            this.ibans.push(cloneDeep(this.current));
                        }
                        this.current = reset();
                        this.dialog = false;
                        this.selectedIndex = undefined;
                    }
                });

            },
            show() {
                this.current = reset();
                this.selectedIndex = undefined;
                const self = this;
                setTimeout(function () {
                    self.dialog = true;
                },100)
            }
        },
        watch:{
            selectedIndex (newValue,oldValue)  {
                console.log(newValue,oldValue);
                if(newValue === undefined){
                    this.current = reset();
                    this.dialog = false;
                }else{
                    this.current = cloneDeep(this.ibans[newValue]);
                    this.dialog = true;
                }
            }
        }
    }
</script>

<style scoped>
</style>