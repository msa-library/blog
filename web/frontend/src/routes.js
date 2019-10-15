//store
import store from './store/index' // your vuex store 

//Site components
import indexPage from './components/Index.vue'

import signInForm from './components/user/SignInForm.vue'
import signUpForm from './components/user/SignUpForm.vue'
import postForm from './components/post/PostForm.vue'
import postPage from './components/post/PostView.vue'
import authorPage from './components/post/PostsByAuthor.vue'
import categoryPage from './components/post/PostsByCategory.vue'
import postAdminPage from './components/post/PostsAdmin.vue'

import categoryForm from './components/category/Form.vue'

const ifAuthenticated = (to, from, next) => {
    if (store.getters.isAuthenticated) {
        next()
        return
    }
    next('/signin')
}

export default [
    {
        path: '/',
        component: indexPage,
        name:'home',
    },
    {
        path: '/signin',
        component: signInForm,
        name:'signin',
    },
    {
        path: '/signup',
        component: signUpForm,
        name:'signup',
    },
    {
        path: '/post/create',
        component: postForm,
        name:'postForm',
        beforeEnter: ifAuthenticated,
    },
    {
        path: '/post/admin',
        component: postAdminPage,
        name:'postAdminPage',
        beforeEnter: ifAuthenticated,
    },

    {
        path: '/author/:slug',
        component: authorPage,
        name:'author',
    },
    {
        path: '/category/create',
        component: categoryForm,
        name:'categoryForm',
        beforeEnter: ifAuthenticated,
    },
    {
        path: '/category/:slug',
        component: categoryPage,
        name:'category',
    },
    {
        path: '/post/:slug',
        component: postPage,
        name:'postPage',
    },

    
];