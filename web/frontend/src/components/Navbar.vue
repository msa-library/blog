<template>
  <div>
    <v-app-bar
      color="primary"
      dense
      dark
    >   
      <v-toolbar-title >
        {{this.$store.state.pageTitle}}
      </v-toolbar-title>

      <div class="flex-grow-1"></div>

      <v-toolbar-items>
        
        <v-btn icon to="/"><v-icon>mdi-home</v-icon></v-btn>
        <v-btn text to="/signin" v-if="this.isGuest">Вход</v-btn>
        <v-btn icon to="/post/create" v-if="!this.isGuest">
          <v-icon>mdi-plus-circle</v-icon>
        </v-btn>
         <v-menu left bottom v-if="!this.isGuest" nudge-bottom="55">
          <template v-slot:activator="{ on }">
            <v-btn icon v-on="on">
              <v-icon>person_outline</v-icon>
            </v-btn>
          </template>

          <v-list>
            <v-list-item to="/post/create">
              <v-list-item-icon><v-icon>note_add</v-icon></v-list-item-icon>
              <v-list-item-title >Новая статья</v-list-item-title>
            </v-list-item>
            <v-list-item to="/post/admin">
              <v-list-item-icon><v-icon>insert_drive_file</v-icon></v-list-item-icon>
              <v-list-item-title >Мои статьи</v-list-item-title>
            </v-list-item>
            <v-divider></v-divider>
            <v-list-item to="/category/create">
              <v-list-item-icon><v-icon>rate_review</v-icon></v-list-item-icon>
              <v-list-item-title >Новая категория</v-list-item-title>
            </v-list-item>
            <v-divider></v-divider>
            <v-list-item>
              <v-list-item-icon><v-icon>person_outline</v-icon></v-list-item-icon>
              <v-list-item-title to="/">Мой профиль</v-list-item-title>
            </v-list-item>
            <v-list-item>
              <v-list-item-icon><v-icon>exit_to_app</v-icon></v-list-item-icon>
              <v-list-item-title @click="logout">Выход</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-toolbar-items>

      

      <!-- <v-btn icon>
        <v-icon>favorite</v-icon>
      </v-btn>

      <v-btn icon>
        <v-icon>search</v-icon>
      </v-btn>

      <v-menu
        left
        bottom
        nudge-bottom="50"
        nudge-left="5"
      >
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>more_vert</v-icon>
          </v-btn>
        </template>

        <v-list>
          <v-list-item to="/signon">
            <v-list-item-title>Вход</v-list-item-title>
          </v-list-item>
          <v-list-item to="/signup">
            <v-list-item-title>Регистация</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu> -->
    </v-app-bar>
  </div>
</template>

<script>
  export default {
    
    data: function () {
      return {
        //isGuest:true
      }
    },
    computed:{
      isGuest:function(){
        if(this.$store.state.authorization!=null){
          return false;
        }
        return true;
      }
    },
    methods:{
      logout:function(){
        
        var self=this;
        this.$store.dispatch('logout').then(function(response){
          
        },function(err){
          
        });
        this.$cookie.delete('Authorization');
        this.$router.push('/');
      }
    },
    mounted:function(){
      
      
      
      /*
      if(this.$store.state.authorization){
        self.isGuest=false
      }else{
        self.isGuest=true
      }
      console.log(self.isGuest)
      */
      
    }
    
  }
</script>