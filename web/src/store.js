import { createStore } from 'vuex'
import axios from 'axios'
import router from './router'

const queryIbanUpdate = `
                    mutation ($id: ID!, $text: String!, $password: String!, $handle: String!, $isPrivate: Boolean!) {
                        ibanUpdate(id: $id, text: $text, password: $password, handle: $handle isPrivate: $isPrivate) {
                            ok,
                            error,
                            iban {id}
                        }
                    }
                `;

const queryIbanCreate = `
                    mutation ($text: String!, $password: String!, $handle: String!, $isPrivate: Boolean!) {
                        ibanNew(text: $text, password: $password, handle: $handle isPrivate: $isPrivate) {
                            ok,
                            error,
                            iban {id}
                        }
                    }
                `;

export default createStore({
    state: {
        token: null,
        error: null,
        isLoaded: false,
        profile: null,
        security: null,
        ibans: [],
        canShow: false,
        logged: false,
    },
    mutations: {
        SET_TOKEN(state, token) {
            localStorage.setItem('user', token);
            axios.defaults.headers.common['Authorization'] = `Bearer ${
                token
            }`;
            state.token = token
        },
        SET_HEADER(state,token) {
            axios.defaults.headers.common['Authorization'] = `Bearer ${
                token
            }`;
            state.token = token;
        },
        SET_IBANS(state,ibans){
            state.ibans = ibans;
        },
        SET_ERROR(state,error){
            state.error = error;
        },
        SET_USER_DATA(state, profile){
            //console.log('SET_USER_DATA');
            state.profile = profile;
            state.logged = true;
        },
        LOGOUT() {
            localStorage.removeItem('user');
            location.reload()
        },
        SET_IS_LOADED(state,isLoaded) {
            state.isLoaded = isLoaded
        },
        CHANGE_PROFILE(state,payload) {
            //console.log(payload);
            state.profile = {
                ...state.profile,
                ...payload
            }
        },
        CHANGE_PASSWORD(state,payload) {
            state.security = {
                password : payload
            }
        },
        SET_SINGLE(state, payload) {
            const {user,iban} = payload;
            state.profile = user;
            state.ibans = iban;
        },

        SET_SHOW_INFO(state, payload) {
            state.canShow = payload;
        },

        RESET_ERROR(state) {
            state.error = null;
        }

    },
    getters: {
        getProfile: (state) => {
            return state.profile;
        }
    },
    actions: {
        resetError({commit}) {
            commit('RESET_ERROR');
        },
        setLoaded({commit},payload) {
            commit('SET_IS_LOADED', payload);
        },
        mapFields({commit,state}, options) {
            const object = {};
            for (let x = 0; x < options.fields.length; x++) {
                const field = [options.fields[x]];
                object[field] = {
                    get() {
                        return state[options.base][field];
                    },
                    set(value) {
                        commit(options.mutation, { [field]: value });
                    }
                };
            }
            return object;
        },

        async ibanDelete({commit}, id) {
            commit('SET_IS_LOADED', false);
            //console.log(id);
            return axios.post('/graph', {
                query: `
                    mutation ($id: ID!) {
                        ibanDelete(id: $id) {
                            ok,
                            error
                        }
                    }
                `,
                variables: {
                    id
                }
            }).then(({data}) => {
                if(data.errors) {
                    alert(data.errors[0].message);
                    return
                }
                return data;
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            });
        },

        fetchIbans({commit}) {
            commit('SET_IS_LOADED', false);
            axios.post('/graph',{
                query: `{
                 getMyIbans{ok,error,iban{id,handle,text,isPrivate}}
                }`,
            }).then(({data}) => {
                //console.log('data');
                //console.log(data);
                if(data.data.getMyIbans.ok) {
                    commit('SET_IBANS', data.data.getMyIbans.iban);
                }
            }).catch((error) => {
                console.log(error)
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            })
        },
        async getUser({commit}) {
            const response = await axios.post('/graph',{
                query: `{
                      getMyProfile {
                        user {id,firstName,lastName,handle,bio,email},
                        ok,
                        error
                      }
                    }`,
            });
            //console.log('response');
            //console.log(response);
            if(response.data.errors){
                router.push('/login');
            }
            if(response.data.data){
                commit('SET_USER_DATA', response.data.data.getMyProfile.user);
            }else{
                router.push('/login');
            }
        },
        fetchProfile({commit}) {
            commit('SET_IS_LOADED', false);
            axios.post('/graph',{
                query: `{
                      getMyProfile {
                        user {id,firstName,lastName,handle,bio,email},
                        ok,
                        error
                      }
                    }`,
            }).then(({data}) => {
                //console.log('data');
                //console.log(data);
                if(data.data.getMyProfile.ok) {
                    commit('SET_USER_DATA', data.data.getMyProfile.user);
                }else{
                    router.push('/login');
                }

            }).catch((error) => {
                console.log(error)
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            })
        },
        checkShowPassword({commit},variables) {
            commit('SET_IS_LOADED', false);
            axios.post('/graph',{
                query: `query ShowInfo($id: ID!,$password: String!) {
                    showInfo(id: $id,password: $password) {
                        ok,
                        error
                    }
                }`,
                variables
            }).then(({data}) => {
                if(data.errors) {
                    alert(data.errors[0].message);
                    return
                }
                commit('SET_SHOW_INFO',data.data.showInfo.ok);
                //console.log(data);
            }).catch((error) => {
                console.log(error)
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            })
        },
        fetchSingleProfile({commit},variables) {
            commit('SET_IS_LOADED', false);
            axios.post('/graph',{
                query: `query GetProfile($username: String!) {
                    getProfile(username: $username) {
                        ok,
                        error,
                        user {
                            visible,
                            handle,
                            firstName,
                            lastName,
                            email,
                            bio
                        },
                        iban {
                            id,
                            text,
                            isPrivate,
                            handle
                        }
                    }
                }`,
                variables
            }).then(({data}) => {
                if(data.errors) {
                    alert(data.errors[0].message);
                    return
                }
                commit('SET_SINGLE',data.data.getProfile);
                //console.log(data);
            }).catch((error) => {
                console.log(error)
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            })
        },
        async ibanUpdate({commit},variables) {
            //console.log(variables);
            commit('SET_IS_LOADED', false);
            let query = "";
            if("id" in variables && variables["id"] !== "") {
                if(variables["isPrivate"] && !("password" in variables)) {
                    variables["password"] = "";
                }
                query = queryIbanUpdate;
            }else{
                query = queryIbanCreate;
            }

            return axios.post('/graph', {
                query,
                variables
            }).then(({data}) => {
                //console.log(data);
                return data;
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            });
        },
        changePassword({commit},credentials) {
            //console.log(credentials);
            commit('SET_IS_LOADED', false);
            //console.log(credentials);
            axios.post('/graph', {
                query: `
                    mutation ($password: String!) {
                        changePassword(password: $password) {
                            ok,
                            error
                        }
                    }
                `,
                variables: {
                    "password" : credentials.password
                }
            }).then(({data}) => {
                console.log(data)
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            });
        },
        changeProfile({commit},variables) {
            //console.log(variables);
            commit('SET_IS_LOADED', false);
            axios.post('/graph', {
                query: `
                    mutation ($bio: String!, $handle: String!) {
                        changeProfile(bio: $bio, handle: $handle) {
                            ok,
                            error
                        }
                    }
                `,
                variables
            }).then(({data}) => {
                if(data.errors){
                    commit('SET_ERROR', data.errors[0].message);
                    return;
                }
                //console.log(data)
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            });
        },
        register({commit}, credentials) {
            commit('SET_IS_LOADED', false);
            axios.post('/graph', {
                query: `
                    mutation ($email: String!, $password: String!, $firstName: String!, $lastName: String!, $handle: String!, $visible: Boolean!) {
                      signUp(email: $email,password: $password, firstName: $firstName, lastName: $lastName, handle: $handle, visible: $visible){
                        ok,
                        error,
                        user{
                          createdAt,
                          updatedAt,
                          handle,
                        },
                      }
                    }`,
                variables: {
                    ...credentials
                },
            }).then(({data}) => {
                //console.log(data);
                //console.log(data.data.signUp.ok);
                if(data.data.signUp.ok){
                    router.push('/login');
                }else{
                    commit('SET_ERROR', data.data.signUp.error);
                }
            }).catch( response => {
                console.log(response);
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            });
        },
        login({commit}, credentials) {
            commit('SET_IS_LOADED', false);
            axios.post('/graph', {
                query: `
                    mutation ($email: String!, $password: String!) {
                      signIn(email: $email,password: $password){
                        ok,
                        error,
                        token,
                      }
                    }`,
                variables: {
                    ...credentials
                },
            }).then(({data}) => {
                if(data.data.signIn.ok){
                    commit('SET_TOKEN', data.data.signIn.token);
                    router.push('/dashboard');
                }else{
                    commit('SET_ERROR', data.data.signIn.error);
                }
            }).catch( response => {
                console.log(response);
            }).finally(() => {
                commit('SET_IS_LOADED', true);
            });
        },
        logout({commit}) {
            commit('LOGOUT')
        },
    }
})