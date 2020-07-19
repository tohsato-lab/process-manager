import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {HttpClientModule} from '@angular/common/http';
import {NgxDropzoneModule} from 'ngx-dropzone';
import {FormsModule} from '@angular/forms';
import {HomeComponent} from './home/home.component';
import {HistoryComponent} from './history/history.component';
import {DatasetComponent} from './dataset/dataset.component';
import {RouterModule, Routes} from '@angular/router';

const routes: Routes = [
    {path: '', component: HomeComponent},
    {path: 'home', component: HomeComponent},
    {path: 'dataset', component: DatasetComponent},
    {path: 'history', component: HistoryComponent},
];

@NgModule({
    declarations: [
        AppComponent,
        HomeComponent,
        HistoryComponent,
        DatasetComponent
    ],

    imports: [
        BrowserModule,
        HttpClientModule,
        NgxDropzoneModule,
        FormsModule,
        RouterModule.forRoot(routes),
    ],

    providers: [],
    bootstrap: [AppComponent]
})


export class AppModule {
}
