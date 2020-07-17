import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {HttpClientModule} from '@angular/common/http';
import {NgxDropzoneModule} from 'ngx-dropzone';
import {FormsModule} from '@angular/forms';

@NgModule({
  declarations: [
    AppComponent
  ],

    imports: [
        BrowserModule,
        HttpClientModule,
        NgxDropzoneModule,
        FormsModule
    ],

  providers: [],
  bootstrap: [AppComponent]
})

export class AppModule {
}
