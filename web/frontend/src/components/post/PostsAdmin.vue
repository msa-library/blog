<template>
    <v-app id="top">
      
      <v-container>
        <BreadcrumbsWidget v-bind:items="breadcrumbs"></BreadcrumbsWidget>
        <v-layout text-xs-center wrap>
          <v-flex xs12>
            <PostListWidget :items="posts"></PostListWidget>
          </v-flex>
          
        </v-layout>
      </v-container>
    </v-app>
    
</template>
<script>
  import BreadcrumbsWidget from '../widgets/Breadcrumbs.vue'
  import PostListWidget from '../widgets/PostList.vue'
  
  export default {
    components: {
      BreadcrumbsWidget,
      PostListWidget,
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

      this.getAuthor(this.$store.state.userId)
      
      
    }
    
  }
</script>
