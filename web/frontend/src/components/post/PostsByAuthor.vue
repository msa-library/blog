<template>
    <v-app id="top">
      
      <v-container>
        <BreadcrumbsWidget v-bind:items="breadcrumbs"></BreadcrumbsWidget>
        <v-layout text-xs-center wrap>
          <v-flex xs12 md9>
            <PostListWidget :items="posts"></PostListWidget>
          </v-flex>
          <v-flex xs12 md3>
            <!-- Автор -->
            <!-- <div class="pa-2">
              <AuthorWidget :item="author"></AuthorWidget>
            </div> -->
            <!-- Категории -->
            <div class="pa-2">
              <CaregoryListWidget v-bind:items="categories"></CaregoryListWidget>
            </div>
            <!-- Авторы -->
            <div class="pa-2">
              <AuthorListWidget v-bind:items="authors" @click="getAuthor"></AuthorListWidget>
            </div>
          </v-flex>
        </v-layout>
      </v-container>
    </v-app>
    
</template>
<script>
  import BreadcrumbsWidget from '../widgets/Breadcrumbs.vue'
  import PostListWidget from '../widgets/PostList.vue'
  import AuthorWidget from '../widgets/Author'
  import AuthorListWidget from '../widgets/AuthorList'
   import CaregoryListWidget from '../widgets/CaregoryList'
  export default {
    components: {
      BreadcrumbsWidget,
      PostListWidget,
      AuthorWidget,
      AuthorListWidget,
      CaregoryListWidget,
    },
    data: function () {
      return {
        breadcrumbs:[],
        author:{
          SrcCover:'',
          FirstName:'',
          LastName:'',
        },
        posts:[],
        authors:[],
        categories:[],
      }
    },
    methods:{
      getFullName:function(){
        return this.author.FirstName+' '+this.author.LastName;
      },
      
      getAuthor:function(slug){
        //Загрузка автора
        var self=this
        this.$store.dispatch('getAuthor',{slug:slug}).then(response=>{
          if(response.body.Author){
            self.author=response.body.Author
            self.breadcrumbs[1].text=self.getFullName()
            self.posts=response.body.Author.Posts
          }
        },err=>{
          console.log("Ошибка загрузки постов: ",err)
        });
      },
    },
    
    
    created:function(){
      
      this.breadcrumbs=[
        {
          text: 'Главная',
          disabled: false,
          to: '/',
        },
        {
          text: '',
          disabled: true,
        },
      ]

      this.getAuthor(this.$route.params.slug)
      
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
