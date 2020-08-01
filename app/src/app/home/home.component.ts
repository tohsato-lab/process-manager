import {Component, OnDestroy, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {webSocket} from 'rxjs/webSocket';
import config from '../../../config';
import {Subscription} from 'rxjs';
import {CommonService} from '../service/commom.service';

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit, OnDestroy {

    public headerTitle = 'プロセス一覧';
    public hiddenUploadPage = true;
    public hiddenAuthPage = true;
    public hiddenCtrlPanel = true;
    public serverAddress = `${config.httpScheme}${location.hostname}:${config.port}`;
    public files: any = [];
    public processList = [];
    public processCtrlData: any = null;

    private subscription: Subscription;

    constructor(
        private http: HttpClient,
        private commonService: CommonService
    ) {
    }

    ngOnInit(): void {
        webSocket(`${config.websocketScheme}${location.hostname}:${config.port}/process_status`).subscribe(
            (message: any) => {
                this.processList = message;
                for (const process of this.processList) {
                    process.Selected = false;
                }
                console.log(this.processList);
            },
            err => {
                console.log(err);
                console.log(`${config.httpScheme}${location.hostname}:${config.port}`);
                this.hiddenAuthPage = false;
            },
            () => console.log('complete')
        );
        this.commonService.onNotifySharedDataChanged(this.headerTitle);
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        this.subscription.unsubscribe();
    }

    public onAddButton(): void {
        this.hiddenUploadPage = false;
    }

    public onCloseUpload(): void {
        this.hiddenUploadPage = true;
        this.files = [];
    }

    public onSelectFiles(event): void {
        console.log(event);
        for (const file of [...event.addedFiles]) {
            this.files.push({data: file, vram: 0.0});
        }
    }

    public onKill(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/kill?id=${id}`
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }

    public onDelete(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/delete?id=${id}`
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }

    public onCloseCtrl(): void {
        this.hiddenCtrlPanel = true;
    }

    public onUpload(): void {
        for (const file of this.files) {
            this.upload(file.data, file.vram);
        }
        this.files = [];
    }

    private upload(file, vram): void {
        console.log('upload');
        const formData = new FormData();
        formData.append('file', file, file.name);
        formData.append('vram', vram);
        this.onCloseUpload();
        this.http.post(
            `${config.httpScheme}${location.hostname}:${config.port}/upload`, formData
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }
}
