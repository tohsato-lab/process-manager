import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {HttpClientModule} from '@angular/common/http';
import {NgxDropzoneModule} from 'ngx-dropzone';
import {FormsModule} from '@angular/forms';
import {HomeComponent} from './home/home.component';
import {HistoryComponent} from './history/history.component';
import {DatasetComponent} from './dataset/dataset.component';
import {CommonService} from './service/commom.service';
import {RouterModule, Routes} from '@angular/router';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {LayoutModule} from '@angular/cdk/layout';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatButtonModule} from '@angular/material/button';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatIconModule} from '@angular/material/icon';
import {MatListModule} from '@angular/material/list';
import {HostsComponent} from './hosts/hosts.component';
import {ChartsModule} from 'ng2-charts';
import {TrashComponent} from './trash/trash.component';

const routes: Routes = [
    {path: '', component: HomeComponent},
    {path: 'home', component: HomeComponent},
    {path: 'dataset', component: DatasetComponent},
    {path: 'hosts', component: HostsComponent},
    {path: 'trash', component: TrashComponent},
];

@NgModule({
    declarations: [
        AppComponent,
        HomeComponent,
        HistoryComponent,
        DatasetComponent,
        HostsComponent,
        TrashComponent,
    ],

    imports: [
        BrowserModule,
        HttpClientModule,
        NgxDropzoneModule,
        FormsModule,
        RouterModule.forRoot(routes, {relativeLinkResolution: 'legacy'}),
        BrowserAnimationsModule,
        LayoutModule,
        MatToolbarModule,
        MatButtonModule,
        MatSidenavModule,
        MatIconModule,
        MatListModule,
        ChartsModule,
    ],

    providers: [CommonService],
    bootstrap: [AppComponent]
})


export class AppModule {
}
