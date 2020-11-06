import {Component, OnDestroy, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {webSocket} from 'rxjs/webSocket';
import config from '../../../config';
import {Subscription} from 'rxjs';
import {CommonService} from '../service/commom.service';

import {SseService} from '../sse.service';

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit, OnDestroy {

    public headerTitle = 'プロセス一覧';
    public hiddenUploadPage = true;
    public hiddenAuthPage = true;
    public serverAddress = `${config.httpScheme}${location.hostname}:${config.port}`;
    public fileInfos: any = [];
    public processList = [];
    public envList: any;
    public processCtrlData: any = null;

    private subscription: Subscription;

    constructor(
        private http: HttpClient,
        private commonService: CommonService,
        private sseService: SseService
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
        this.sseService
            .getServerSentEvent(`${config.httpScheme}${location.hostname}:${config.port}/gpu_status`)
            .subscribe(data => console.log(data.data));
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        this.subscription.unsubscribe();
    }

    public onOpenExplorer(id) {
        window.location.href = "http://localhost:5983/programs/" + id;
    }

    public onAddButton(): void {
        this.hiddenUploadPage = false;
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/env_info`
        ).subscribe(value => {
            this.envList = value;
        }, error => {
            console.log(error);
        });
    }

    public onCloseUpload(): void {
        this.hiddenUploadPage = true;
        this.fileInfos = [];
    }

    public onSelectFiles(event): void {
        console.log(event);
        for (const file of [...event.addedFiles]) {
            this.fileInfos.push({file: file, vram: 0.0, env: 'base', target: 'main.py'});
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

    public onUpload(): void {
        for (const fileInfo of this.fileInfos) {
            this.upload(fileInfo);
        }
        this.fileInfos = [];
    }

    private upload(info): void {
        console.log('upload');
        const formData = new FormData();
        formData.append('file', info.file, info.file.name);
        formData.append('vram', info.vram);
        formData.append('env', info.env);
        formData.append('target', info.target);
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
