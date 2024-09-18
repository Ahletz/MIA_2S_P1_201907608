import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ComandoComponent } from './comando/comando.component';

import { FormsModule } from '@angular/forms';  // Importa FormsModule para usar ngModel
import { HttpClientModule } from '@angular/common/http';  // Para hacer peticiones HTTP

@NgModule({
  declarations: [
    AppComponent,
    ComandoComponent
  ],
  imports: [
    BrowserModule,
    FormsModule, 
    HttpClientModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
