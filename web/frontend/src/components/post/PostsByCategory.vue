<template>
    <v-app id="top">
      
      <v-container>
        <BreadcrumbsWidget v-bind:items="breadcrumbs"></BreadcrumbsWidget>
        <v-layout text-xs-center wrap>
          <v-flex xs12 md9>
            <PostListWidget v-bind:items="posts"></PostListWidget>
          </v-flex>
          <v-flex xs12 md3>
            <!-- Категории -->
            <div class="pa-2">
              <CaregoryListWidget v-bind:items="categories" @click="getCategory"></CaregoryListWidget>
            </div>

            <!-- Авторы -->
            <div class="pa-2">
              <AuthorListWidget v-bind:items="authors"></AuthorListWidget>
            </div>
          </v-flex>
        </v-layout>
      </v-container>
    </v-app>
    
</template>
<script>
  import BreadcrumbsWidget from '../widgets/Breadcrumbs'
  import PostListWidget from '../widgets/PostList'
  import CaregoryListWidget from '../widgets/CaregoryList'
  import AuthorListWidget from '../widgets/AuthorList'
  export default {
    components: {
      BreadcrumbsWidget,
      PostListWidget,
      CaregoryListWidget,
      AuthorListWidget,
    },
    data: function () {
      return {
        breadcrumbs:[],
        category:null,
        categories:[],
        authors:[],
        posts:[],
      }
    },
    methods:{
        getCategory:function(slug){
          var self=this
          //Загрузка текущей категории
          this.$store.dispatch('getPostCategory',{slug:slug}).then(response=>{
            if(response.body.Category){
              self.category=response.body.Category
              self.posts=response.body.Category.Posts
              self.breadcrumbs=[
                {
                  text: 'Главная',
                  disabled: false,
                  to: '/',
                },
                {
                  text: response.body.Category.Name,
                  disabled: true,
                },
              ];
            }
          },err=>{
            console.log("Ошибка загрузки категорий: ",err)
          });
        }
    },
    created:function(){
      this.getCategory(this.$route.params.slug)


      var self=this
      //Загрузка списка всех категорий
      this.$store.dispatch('getCategories',{}).then(response=>{
        if(response.body.Categories){
          self.categories=response.body.Categories
        }
      },err=>{
        console.log("Ошибка загрузки категорий: ",err)
      });

      //Загрузка авторов
      this.$store.dispatch('getAuthors',{}).then(response=>{
        
        if(response.body.Authors){
          self.authors=response.body.Authors
          
        }
      },err=>{
        console.log("Ошибка загрузки авторов: ",err)
      });
    }
    
  }
</script>
