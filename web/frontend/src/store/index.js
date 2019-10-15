import Vue from 'vue'
import Vuex from 'vuex'
Vue.use(Vuex)

// VueCookie
import VueCookie from 'vue-cookie';
Vue.use(VueCookie);

import VueResource from 'vue-resource'
Vue.use(VueResource)

Vue.http.options.root = process.env.VUE_APP_API_ENDPOINT
//Vue.http.options.credentials = true;
//Vue.http.headers.common['Access-Control-Allow-Origin'] = '*';
Vue.http.interceptors.push(function(request, next) {
  request.headers.set('Authorization', VueCookie.get('Authorization'))
  next()
});


const store= new Vuex.Store({
  
  state:{
    pageTitle:'MSA Blog',
    authorization:null,
    userId:null,
    role:null,
  },
  getters:{
    isAuthenticated: state => state.authorization?true:false,
    isGuest: state => !state.authorization?true:false,
  },
  mutations:{
    set(state,{type,data}){
      state[type]=data
    },
    push(state,{type,data}){
      state[type].push(data)
    }
  },
  actions:{
    //Вход
    signin({commit},data){
      const url='user/signin';
      return new Promise((resolve, reject) => {
        Vue.http.post(url,data,{emulateJSON: true}).then(function(response){
          
          var token=response.headers.get('Authorization')
          commit('set',{type:'authorization',data:token});
          commit('set',{type:'userId',data:response.body.Slug});
          commit('set',{type:'role',data:response.body.Role});
          
          resolve(response)
        },function(err){
          reject(err)
        })  
      })
    },
    
    //Регистрация
    signup({commit},data){
      const url='user/signup';
      return new Promise((resolve, reject) => {
        
        Vue.http.post(url,data,{emulateJSON: true}).then(function(response){
          
          var token=response.headers.get('Authorization')
          commit('set',{type:'authorization',data:token});
          
          resolve(response)
        },function(err){
          reject(err)
        })  

      })
    },

    //Выход
    logout({commit},data){
      const url='user/logout';
      return new Promise((resolve, reject) => {
        Vue.http.get(url).then(function(response){

          //Logout
          commit('set',{type:'authorization',data:null});
          delete Vue.http.headers.common['Authorization'];
          
          resolve(response);
          
        },function(err){
          commit('set',{type:'authorization',data:null});
          reject(err);
        })  
      });
    },

    //Создание статьи
    createPost({commit},data){
      const url='post';
      return new Promise((resolve, reject) => {
        Vue.http.post(url,data,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Список статей
    getPosts({commit},data){
      const url='post';
      return new Promise((resolve, reject) => {
        Vue.http.get(url,data,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Статья по ID
    getPost({commit},data){
      const url='post/'+data.slug;
      return new Promise((resolve, reject) => {
        Vue.http.get(url,data,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Список Авторов
    getAuthors({commit},data){
      const url='author';
      return new Promise((resolve, reject) => {
        Vue.http.get(url,data,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Автор по ID
    getAuthor({commit},data){
      const url='author/'+data.slug;
      return new Promise((resolve, reject) => {
        Vue.http.get(url,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Список категорий
    getCategories({commit},data){
      const url='category';
      return new Promise((resolve, reject) => {
        Vue.http.get(url,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    // PostCategory по ID
    getPostCategory({commit},data){
      const url='post/category/'+data.slug;
      return new Promise((resolve, reject) => {
        Vue.http.get(url,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    // Категории по ID
    getCategory({commit},data){
      const url='category/'+data.slug;
      return new Promise((resolve, reject) => {
        Vue.http.get(url,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Создание категории
    createCategory({commit},data){
      const url='category';
      return new Promise((resolve, reject) => {
        Vue.http.post(url,data,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },

    //Создание коментария
    createComment({commit},data){
      const url='comment';
      return new Promise((resolve, reject) => {
        Vue.http.post(url,data,{
          emulateJSON: true
        }).then(function(response){
          resolve(response);
        },function(err){
          reject(err);
        })  
      });
    },
  }
})

export default store